package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	driver "github.com/go-sql-driver/mysql"

	"sentence.valc.cc/backend/internal/config"
	"sentence.valc.cc/backend/internal/domain"
	apihttp "sentence.valc.cc/backend/internal/infra/http"
	"sentence.valc.cc/backend/internal/infra/provider/ninja"
	mysqlrepo "sentence.valc.cc/backend/internal/infra/repository/mysql"
	"sentence.valc.cc/backend/internal/usecase"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := openDB(cfg.MySQLDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	repo := mysqlrepo.NewQuoteRepository(db)
	provider := ninja.New(cfg.NinjaAPIKey)

	qotd := usecase.NewGetQuoteOfTheDay(repo, provider, domain.RealClock{})
	react := usecase.NewReactToQuote(repo)

	handler := apihttp.NewHandler(qotd, react)
	router := apihttp.CORS(cfg.AllowOrigins, apihttp.RateLimit(60, apihttp.NewRouter(handler)))

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("http listening", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server", "err", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}

// openDB connects to MySQL and retries until the server is reachable.
func openDB(dsn string) (*sql.DB, error) {
	// Validate DSN early for a clearer error.
	if _, err := driver.ParseDSN(dsn); err != nil {
		return nil, err
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	var pingErr error
	for i := 0; i < 30; i++ {
		if pingErr = db.Ping(); pingErr == nil {
			slog.Info("connected to mysql")
			return db, nil
		}
		slog.Warn("waiting for mysql", "attempt", i+1, "err", pingErr)
		time.Sleep(2 * time.Second)
	}
	return nil, pingErr
}
