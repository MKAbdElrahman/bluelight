package main

import (
	"context"
	"database/sql"
	"errors"
	"expvar"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"bluelight.mkcodedev.com/src/api/handlers/middleware"
	"bluelight.mkcodedev.com/src/api/router"
	"bluelight.mkcodedev.com/src/infrastructure/mailer"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const version = "0.0.3"

type serverConfig struct {
	port            int
	env             string
	version         bool
	readTimeout     time.Duration
	writeTimeout    time.Duration
	idleTimeout     time.Duration
	shutdownTimeout time.Duration
}

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}
type cors struct {
	trustedOrigins []string
}
type smtp struct {
	host   string
	sender string
}

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
		server  serverConfig
		db      dbConfig
		limiter middleware.RateLimiterConfig
		smtp    smtp
		cors    cors
	}

	flag.IntVar(&cfg.server.port, "port", 3000, "API server port")
	flag.StringVar(&cfg.server.env, "env", "development", "Environment (development|staging|production)")
	flag.DurationVar(&cfg.server.readTimeout, "server-read-timeout", 5*time.Second, "Maximum duration for reading the entire request, including the body.")
	flag.DurationVar(&cfg.server.writeTimeout, "server-write-timeout", 10*time.Second, "Maximum duration for writing the response, including the body.")
	flag.DurationVar(&cfg.server.idleTimeout, "server-idle-timeout", 120*time.Second, "Maximum amount of time to wait for the next request when keep-alives are enabled.")
	flag.DurationVar(&cfg.server.shutdownTimeout, "server-shutdown-timeout", 30*time.Second, "Maximum duration to wait for active connections to close during server shutdown.")
	flag.BoolVar(&cfg.server.version, "version", false, "Show API version")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.RPS, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-base-url", "https://sandbox.api.mailtrap.io", "SMTP base url")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "no-reply@bluelight.mkcodedev.net", "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	flag.Parse()

	if cfg.server.version {
		fmt.Printf("API Version: %s\n", version)
		os.Exit(0)
	}

	// This is just for test!!
	mailtrapURL := "https://sandbox.api.mailtrap.io/api/send/3109121"
	apiKey := "245e889dc812ac10b67d19e4801ba7e7"
	senderEmail := cfg.smtp.sender

	mailerClient := mailer.NewMailer(mailer.Config{
		URL:    mailtrapURL,
		APIKey: apiKey,
		Sender: mailer.EmailAddress{
			Email: senderEmail,
			Name:  "Bluelight",
		},
	})

	// POSTGRESQL
	db, err := openDB(cfg.db)
	if err != nil {
		logger.Error("databasse connection failed", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("databasse connection pool established")

	// ROUTER
	backgroundRoutinesWaitGroup := &sync.WaitGroup{}
	router := router.NewRouter(router.RouterConfig{
		TrustedOrigins:      cfg.cors.trustedOrigins,
		Logger:              logger,
		API_Environment:     cfg.server.env,
		API_Version:         version,
		DB:                  db,
		Mailer:              mailerClient,
		LimiterConfig:       cfg.limiter,
		BackgroundWaitGroup: backgroundRoutinesWaitGroup,
	})

	expvar.NewString("version").Set(version)
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	// SERVER
	err = serve(backgroundRoutinesWaitGroup, logger, router, cfg.server)
	if err != nil {
		logger.Error("server failed to gracefully shutdown", "error", err.Error())
		os.Exit(1)
	}

}

func serve(backgroundRoutinesWaitGroup *sync.WaitGroup, logger *slog.Logger, r http.Handler, cfg serverConfig) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      r,
		ReadTimeout:  cfg.readTimeout,
		IdleTimeout:  cfg.idleTimeout,
		WriteTimeout: cfg.writeTimeout,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		logger.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), cfg.shutdownTimeout)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		logger.Info("completing background tasks", "addr", srv.Addr)
		backgroundRoutinesWaitGroup.Wait()
		shutdownError <- nil
	}()

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	err := srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	logger.Info("stopped server", "addr", srv.Addr)

	return nil
}

func openDB(cfg dbConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.maxOpenConns)
	db.SetMaxIdleConns(cfg.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
