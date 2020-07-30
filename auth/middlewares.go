package auth

import (
	"context"
	"net/http"
)

// LoadContextMiddleware puts information about current user into request context.
// This middleware is required and should be connected to mux of a route.
func LoadContextMiddleware(as Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(r.Header["Authorization"]) == 1 {
				tokenStr := r.Header["Authorization"][0]

				user, err := as.FindUserByToken(r.Context(), tokenStr)

				if err == nil {
					r = r.WithContext(
						context.WithValue(
							r.Context(),
							CtxKeyAuthUser,
							user,
						),
					)
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
