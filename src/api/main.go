package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"bluelight.mkcodedev.com/src/api/handlers"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const version = "0.0.2"

func main() {
	// LOGGING
	logHandler := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
	})
	logger := slog.New(logHandler)

	// CONFIGURATION
	err := godotenv.Load("resources/.env")
	if err != nil {
		logger.Error("couldn't load .env file")
	}

	var cfg struct {
		port    int
		env     string
		version bool
		db      struct {
			dsn          string
			maxOpenConns int
			maxIdleConns int
			maxIdleTime  time.Duration
		}
	}

	flag.IntVar(&cfg.port, "port", 3000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("BLUELIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")
	flag.BoolVar(&cfg.version, "version", false, "Show API version")

	flag.Parse()

	if cfg.version {
		fmt.Printf("API Version: %s\n", version)
		os.Exit(0)
	}

	// POSTGRESQL
	db, err := openDB(cfg.db.dsn, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Error("databasse connection failed", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("databasse connection pool established")

	// ROUTER
	mux := handlers.NewRouter(handlers.RouterConfig{
		Logger:          logger,
		API_Environment: cfg.env,
		API_Version:     version,
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
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string, maxOpenConns int, maxIdleConns int, maxIdleTime time.Duration) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
