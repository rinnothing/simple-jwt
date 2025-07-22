package api

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/rinnothing/simple-jwt/internal/api/authapi"
	"github.com/rinnothing/simple-jwt/internal/api/schema"
	"github.com/rinnothing/simple-jwt/internal/config"
	"github.com/rinnothing/simple-jwt/internal/repository/postgres"
	"github.com/rinnothing/simple-jwt/internal/service/auth"
	storage "github.com/rinnothing/simple-jwt/internal/service/secure_storage"
	webhook "github.com/rinnothing/simple-jwt/internal/service/webhook_caller"
	migrations "github.com/rinnothing/simple-jwt/postgres"
)

type Server struct {
	cancel context.CancelFunc
}

func (s *Server) Run(cfg config.Config, logger *zap.Logger) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	s.cancel = cancel

	dbPool, err := pgxpool.New(ctx, cfg.Postgres.URL)
	if err != nil {
		logger.Error("cannot connect to database", zap.Error(err))
		return err
	}
	defer dbPool.Close()

	migrations.SetupPostgres(dbPool, logger)

	repo := postgres.NewRepo(cfg.Postgres, dbPool, logger)

	webhook := webhook.NewService(cfg.Webhook, logger)

	storage := storage.NewService(repo, logger)

	auth, err := auth.NewService(&cfg.Auth, repo, webhook, logger)
	if err != nil {
		logger.Error("cannot create auth service", zap.Error(err))
		return err
	}

	serviceAPI := authapi.NewAPI(auth, storage, logger)

	e := echo.New()
	e.Use(echomiddleware.Recover())

	outputs, outputsSimple, err := openEchoOutputs(e, cfg)
	if err != nil {
		logger.Error("cannot open echo outputs", zap.Error(err))
		return err
	}
	defer closeEchoOutputs(outputs)

	e.Use(echomiddleware.LoggerWithConfig(echomiddleware.LoggerConfig{
		Output: io.MultiWriter(
			outputsSimple...,
		),
	}))

	e.IPExtractor = echo.ExtractIPDirect()

	schema.RegisterHandlers(e, serviceAPI)

	go func() {
		if err := e.Start(net.JoinHostPort("0.0.0.0", cfg.Port)); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("server died", zap.Error(err))
		}
	}()

	<-ctx.Done()

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer stopCancel()

	if err := e.Shutdown(stopCtx); err != nil {
		logger.Fatal("server shutdown failed", zap.Error(err))
		return err
	}

	logger.Info("server stopped")
	return nil
}

func (s *Server) Stop() {
	s.cancel()
}

func openEchoOutputs(e *echo.Echo, cfg config.Config) ([]io.WriteCloser, []io.Writer, error) {
	outputs := make([]io.WriteCloser, 0)
	outputsSimple := make([]io.Writer, 0)
	for _, path := range cfg.Logger.EchoOutputs {
		if path == "stdout" {
			outputs = append(outputs, os.Stdout)
			outputsSimple = append(outputsSimple, os.Stdout)
			continue
		} else if path == "stderr" {
			outputs = append(outputs, os.Stderr)
			outputsSimple = append(outputsSimple, os.Stderr)
			continue
		}

		output, err := os.OpenFile(path, os.O_RDWR, 0777)
		if err != nil {
			return nil, nil, err
		}
		outputs = append(outputs, output)
		outputsSimple = append(outputsSimple, output)
	}
	return outputs, outputsSimple, nil
}

func closeEchoOutputs(outputs []io.WriteCloser) error {
	for _, output := range outputs {
		err := output.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
