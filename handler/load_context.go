package handler

import (
	"context"
	"net/http"

	"bitbucket.org/iwlab-standuply/slackteams-api/shared"
)

type CtxKey string

const (
	CtxKeyRequestID CtxKey = "requestId"
)

func LoadContextMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := shared.RandStringBytesMaskImprSrcUnsafe(24)

			r = r.WithContext(
				context.WithValue(
					r.Context(),
					CtxKeyRequestID,
					requestID,
				),
			)

			next.ServeHTTP(w, r)
		})
	}
}
