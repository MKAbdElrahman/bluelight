package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"bluelight.mkcodedev.com/src/api/handlers"
	"github.com/charmbracelet/log"
)

const version = "0.0.1"

func main() {

	// CONFIGURATION
	var cfg struct {
		port int
		env  string
	}

	flag.IntVar(&cfg.port, "port", 3000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// LOGGING
	logHandler := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
	})
	logger := slog.New(logHandler)

	// ROUTER
	mux := handlers.NewRouter(handlers.MuxConfig{
		Logger:      logger,
		Environment: cfg.env,
		Version:     version,
	})
	// SERVER
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
