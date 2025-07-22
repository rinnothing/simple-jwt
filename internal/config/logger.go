package config

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type LoggerConfig struct {
	Env         string   `yaml:"env"`
	OutputPaths []string `yaml:"output_paths,omitempty"`
}

func ConfigureLogger(cfg LoggerConfig) (zap.Config, error) {
	var config zap.Config
	if cfg.Env == "dev" {
		config = zap.NewDevelopmentConfig()
	} else if cfg.Env == "prod" {
		config = zap.NewProductionConfig()
	} else {
		fmt.Fprintf(os.Stderr, "no such logger environment: %s", cfg.Env)
		return zap.Config{}, fmt.Errorf("no such logger environment: %s", cfg.Env)
	}

	if cfg.OutputPaths != nil {
		config.OutputPaths = cfg.OutputPaths
	}
	for _, path := range config.OutputPaths {
		if path != "stdout" && path != "stderr" {
			err := os.MkdirAll(filepath.Dir(path), 0777)
			if err != nil {
				return zap.Config{}, fmt.Errorf("can't create dirs: %w", err)
			}
			_, err = os.Create(path)
			if err != nil {
				return zap.Config{}, fmt.Errorf("can't create log file: %w", err)
			}
		}
	}

	return config, nil
}
