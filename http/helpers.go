package http

import (
	"context"
	"errors"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type ctxKey int

const (
	ctxLoggerKey ctxKey = iota
	ctxPageKey
	ctxSizeKey
)

func loggerMiddleware(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			reqLog := log.With(zap.String("req-id", reqID))
			reqLog.Debug("Serving Request")
			r = r.WithContext(context.WithValue(r.Context(), ctxLoggerKey, reqLog))

			next.ServeHTTP(w, r)
		})
	}
}

func getLog(r *http.Request) (*zap.Logger, error) {
	log, ok := r.Context().Value(ctxLoggerKey).(*zap.Logger)
	if !ok {
		return nil, errors.New("Logger not present")
	}
	return log, nil
}

func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			page, size int64
			err        error
		)
		page, size = 1, 25
		pageString := chi.URLParam(r, "page")
		sizeString := chi.URLParam(r, "size")

		if pageString != "" {
			page, err = strconv.ParseInt(pageString, 10, 63)
			if err != nil {
				errorWriter(w, r, http.StatusBadRequest, err)
				return
			}
		}
		if sizeString != "" {
			size, err = strconv.ParseInt(sizeString, 10, 63)
			if err != nil {
				errorWriter(w, r, http.StatusBadRequest, err)
				return
			}
			if size < 1 {
				size = 25
			}
		}
		r = r.WithContext(context.WithValue(r.Context(), ctxPageKey, page))
		r = r.WithContext(context.WithValue(r.Context(), ctxSizeKey, size))

		next.ServeHTTP(w, r)
	})
}

type errorResponse struct {
	Error     string `json:"error,omitempty"`
	RequestID string `json:"req_id,omitempty"`
}

func errorWriter(w http.ResponseWriter, r *http.Request, status int, err error) {
	var resp = &errorResponse{
		Error:     err.Error(),
		RequestID: middleware.GetReqID(r.Context()),
	}
	w.WriteHeader(status)
	render.JSON(w, r, resp)
	return
}

func errorStringWriter(w http.ResponseWriter, r *http.Request, status int, err string) {
	errorWriter(w, r, status, errors.New(err))
}
