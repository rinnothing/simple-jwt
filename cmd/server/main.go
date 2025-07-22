package main

import (
	"fmt"
	"os"

	"github.com/rinnothing/simple-jwt/internal/api"
	"github.com/rinnothing/simple-jwt/internal/config"
	"go.uber.org/zap"
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

	server := api.Server{}
	if err := server.Run(cfg, logger); err != nil {
		logger.Fatal("server stopped with error", zap.Error(err))
	}
}
