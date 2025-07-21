package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Auth     AuthConfig     `yaml:"auth"`
	Postgres PostgresConfig `yaml:"postgres"`
	Webhook  WebhookConfig  `yaml:"webhook"`
}

func GetConfig() (Config, error) {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		return Config{}, fmt.Errorf("can't open config file: %w", err)
	}
	defer file.Close()

	var config Config
	d := yaml.NewDecoder(file)

	err = d.Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("can't unmarshal config file: %w", err)
	}

	err = env.Parse(&config)
	if err != nil {
		return Config{}, fmt.Errorf("can't parse env variables: %w", err)
	}

	config.Postgres.MakeURL()

	return config, nil
}
