package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

const (
	requestIDKey int = 0
)

var (
	listenAddr string
	healthy    int32
)

func main() {
	flag.StringVar(&listenAddr, "listen-addr", ":5000", "server listen address")
	flag.Parse()

	// set logger
	logger := newLogger("dragonfly", "./dragonfly.log", true)
	logger.Info("Server is starting...")

	// http
	router := http.NewServeMux()
	router.Handle("/", index())
	router.Handle("/healthz", healthz())

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	server := &http.Server{
		Addr:    listenAddr,
		Handler: Tracing(nextRequestID)(Logging(logger)(router)),
		// ErrorLog:     logger,      // use log =>  logger := log.New(os.Stdout, "dragonfly: ", log.LstdFlags)
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Info("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cannel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cannel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Could not gracefully shutdown the server: ", err)
		}
		close(done)
	}()

	logger.Info("Server is ready to handle requests at", listenAddr)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Could not listen on : ", listenAddr, err)
	}

	<-done
	logger.Info("Server stopped")
}
