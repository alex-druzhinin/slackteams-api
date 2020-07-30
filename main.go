package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"bitbucket.org/iwlab-standuply/slackteams-api/database"
	"bitbucket.org/iwlab-standuply/slackteams-api/database/mongodb"

	"bitbucket.org/iwlab-standuply/slackteams-api/auth"
	"bitbucket.org/iwlab-standuply/slackteams-api/config"
	"bitbucket.org/iwlab-standuply/slackteams-api/handler"
	"bitbucket.org/iwlab-standuply/slackteams-api/logger"

	log "github.com/sirupsen/logrus"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Headers", "Cookie, Content-Type, X-Auth-Token, X-Language, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.WriteHeader(200)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	conf, err := config.LoadConfig()

	logger.InitLogger(conf)

	if err != nil {
		log.WithError(err).Fatal(`Failed to load config`)
	}

	authService, err := auth.NewAuthService(auth.Config{
		UserRepository: database.NewLocalAuthRepository([]config.User{
			conf.BotUser,
			conf.MeteorUser,
		}),
	})

	if err != nil {
		log.WithError(err).Fatal(`Failed to init AuthService`)
	}

	// Register handlers to routes.
	mux := http.NewServeMux()
	mux.Handle("/", handler.Empty{})

	h := handler.AllAuthorizations{
		Repo: mongodb.NewSlackBotAuthorizationsRepository(conf.MongoDB.URI),
	}
	hndlr := handler.LoadContextMiddleware()(
		auth.LoadContextMiddleware(authService)(
			CorsMiddleware(h),
		),
	)

	mux.Handle("/allAuthorizations/", hndlr)
	mux.Handle("/allAuthorizations", hndlr) // Register without a trailing slash to avoid redirect.

	var (
		readHeaderTimeout = 1 * time.Second
		writeTimeout      = 120 * time.Second
		idleTimeout       = 90 * time.Second
		maxHeaderBytes    = http.DefaultMaxHeaderBytes
	)

	// Configure the HTTP server.
	s := &http.Server{
		Addr:              conf.Addr,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	go func() {
		// Begin listening for requests.
		log.Printf("Listening for requests on %s", s.Addr)

		if err = s.ListenAndServe(); err != nil {
			log.WithError(err).Fatal("ListenAndServe failed")
		}
	}()

	<-stop

	log.Println("Shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // releases resources if s.Shutdown completes before timeout elapses

	if err := s.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Server stopped with errors.")
	} else {
		log.Println("Server gracefully stopped.")
	}
}
