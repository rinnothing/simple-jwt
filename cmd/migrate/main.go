package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/rinnothing/simple-jwt/internal/config"
	"github.com/rinnothing/simple-jwt/postgres"
)

func main() {
	cfg, err := config.GetConfig("config/config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't read config: %s", err.Error())
		panic(err)
	}

	loggerCfg, err := config.ConfigureLogger(cfg.Logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't read logger configuration: %s", err.Error())
		panic(err)
	}

	logger, err := loggerCfg.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't start logger: %s", err.Error())
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.Postgres.URL)
	if err != nil {
		logger.Error("cannot connect to database", zap.Error(err))
		panic(err)
	}
	defer dbPool.Close()

	logger.Info("start migrations")

	postgres.SetupPostgres(dbPool, logger)

	logger.Info("successfully done migrations")
}
