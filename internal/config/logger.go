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
	EchoOutputs []string `yaml:"echo_outputs,omitempty"`
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
		err := tryCreateOut(path)
		if err != nil {
			return zap.Config{}, err
		}
	}
	for _, path := range cfg.EchoOutputs {
		err := tryCreateOut(path)
		if err != nil {
			return zap.Config{}, err
		}
	}

	return config, nil
}

func tryCreateOut(path string) error {
	if path == "stdout" || path == "stderr" {
		return nil
	}

	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return fmt.Errorf("can't create dirs: %w", err)
	}
	_, err = os.Create(path)
	if err != nil {
		return fmt.Errorf("can't create log file: %w", err)
	}

	return nil
}
