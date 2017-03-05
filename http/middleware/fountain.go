package middleware

import (
	"context"
	"iris.arke.works/forum/http/ctxkeys"
	"iris.arke.works/forum/snowflakes"
	"net/http"
)

func FountainMiddleware(fountain snowflakes.Fountain) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), ctxkeys.CtxFountainKey, fountain))

			next.ServeHTTP(w, r)
		})
	}
}
