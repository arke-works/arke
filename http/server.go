package http

import (
	"context"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Start will listen on a given TCP Address/Port and log to a zap.Logger instance.
//
// It returns a struct{} and an error channel, the later will return any errors
// caused by the http server itself and the former is used to signal shutdown
func Start(addr *net.TCPAddr, log *zap.Logger) (chan<- struct{}, <-chan error) {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer,
		middleware.RequestID,
		middleware.CloseNotify,
		middleware.DefaultCompress,
		middleware.RedirectSlashes)

	router.Use(loggerMiddleware(log))
	router.Use(paginate)

	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/:resource/:snowflake", getHandler)
		r.Options("/:resource", optionHandler)
	})

	server := &http.Server{
		Addr:    addr.String(),
		Handler: router,
	}

	// This line prevents a theoretical snowflake collision
	time.Sleep(time.Second)

	shutdownChan := make(chan struct{})
	errorChan := make(chan error)

	go func() {
		<-shutdownChan
		log.Warn("Server shutdown requested")
		server.Shutdown(context.Background())
	}()

	go func() {
		errorChan <- server.ListenAndServe()
		close(errorChan)
	}()

	return shutdownChan, errorChan
}

// StartBlocking works like Start but instead of returning it sets up
// a signal loop and gracefully shuts down the http server if an OS interrupt
// is received.
func StartBlocking(addr *net.TCPAddr, log *zap.Logger) error {
	log.Info("Starting HTTP Server")
	shutdownChan, errorChan := Start(addr, log)

	log.Debug("Setting up shutdown listener")
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt)
	go func() {
		for _ = range osChan {
			shutdownChan <- struct{}{}
			return
		}
	}()

	log.Info("HTTP Server started")
	return <-errorChan
}
