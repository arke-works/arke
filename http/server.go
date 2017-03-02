package http

import (
	"context"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"go.uber.org/zap"
	"iris.arke.works/forum/http/handlers"
	amiddleware "iris.arke.works/forum/http/middleware"
	"iris.arke.works/forum/snowflakes"
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
	fountain := &snowflakes.Generator{
		InstanceID: 1,
		StartTime:  time.Date(2017, 02, 18, 17, 03, 33, 0, time.UTC).Unix(),
	}

	router.Use(middleware.Recoverer,
		middleware.RequestID,
		middleware.CloseNotify,
		middleware.DefaultCompress,
		middleware.RedirectSlashes)

	router.Use(amiddleware.LoggerMiddleware(log))
	router.Use(amiddleware.FountainMiddleware(fountain))
	router.Use(amiddleware.PageMiddleware)

	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/:resource/:snowflake", handlers.GetHandler)
		r.Get("/:resource", handlers.GetHandler)
		r.Head("/:resource/:snowflake", handlers.GetHandler)
		r.Head("/:resource", handlers.GetHandler)
		r.Options("/:resource", handlers.OptionHandler)
		r.Options("/:resource/:unused", handlers.OptionHandler)
		r.Post("/:resource", handlers.PostHandler)
		r.Post("/:resource/:unused", handlers.DenyHandler)
		r.Delete("/:resource/", handlers.DenyHandler)
		r.Delete("/:resource/:snowflake", handlers.DeleteHandler)
		r.Put("/:resource/", handlers.DenyHandler)
		r.Put("/:resource/:snowflake", handlers.DenyHandler)
		r.Patch("/:resource/", handlers.DenyHandler)
		r.Patch("/:resource/:swowflake", handlers.DenyHandler)
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
