package middleware

import (
	"context"
	"github.com/pressly/chi/middleware"
	"go.uber.org/zap"
	"iris.arke.works/forum/http/ctxkeys"
	"net/http"
)

func LoggerMiddleware(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			reqLog := log.With(zap.String("req-id", reqID))
			reqLog.Debug("Serving Request")
			r = r.WithContext(context.WithValue(r.Context(), ctxkeys.CtxLoggerKey, reqLog))

			next.ServeHTTP(w, r)
		})
	}
}
