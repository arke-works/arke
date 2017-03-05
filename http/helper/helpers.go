package helper

import (
	"errors"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
	"go.uber.org/zap"
	"iris.arke.works/forum/http/ctxkeys"
	"net/http"
)

var GetLog = func(r *http.Request) (*zap.Logger, error) {
	log, ok := r.Context().Value(ctxkeys.CtxLoggerKey).(*zap.Logger)
	if !ok {
		return nil, errors.New("Logger not present")
	}
	return log, nil
}

func SetupTestLog() {
	GetLog = func(_ *http.Request) (*zap.Logger, error) {
		return zap.NewDevelopment()
	}
}

type errorResponse struct {
	Error     string `json:"error,omitempty"`
	RequestID string `json:"req_id,omitempty"`
}

func ErrorWriter(w http.ResponseWriter, r *http.Request, status int, err error) {
	var resp = &errorResponse{
		Error:     err.Error(),
		RequestID: middleware.GetReqID(r.Context()),
	}
	w.WriteHeader(status)
	render.JSON(w, r, resp)
	return
}

func ErrorStringWriter(w http.ResponseWriter, r *http.Request, status int, err string) {
	ErrorWriter(w, r, status, errors.New(err))
}
