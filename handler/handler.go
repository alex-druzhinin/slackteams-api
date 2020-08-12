// Package handler defines the HTTP handlers for this API.
package handler

import (
	"net/url"
	"bytes"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func respond(w http.ResponseWriter, body []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, err := w.Write(body)
	if err != nil {
		log.WithError(err).Debug("Write in respond failed")
	}
}

func isSupported(method string) bool {
	return method == "GET"
}

func errorJSON(msg string) []byte {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, `{"error": "%s"}`, msg)
	return buf.Bytes()
}

type request struct {
	queries  url.Values
	authRole AuthRole
}

type AuthRole string

const (
	RoleBot    AuthRole = "standuply-bot"
	RoleMeteor AuthRole = "meteor"
)
