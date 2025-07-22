package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Auth     AuthConfig     `yaml:"auth"`
	Postgres PostgresConfig `yaml:"postgres"`
	Webhook  WebhookConfig  `yaml:"webhook"`
	Logger   LoggerConfig   `yaml:"logger"`
	Port     string         `yaml:"port"`
}

func GetConfig(path string) (Config, error) {
	file, err := os.Open(path)
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

	config.Postgres.MakeURL()

	return config, nil
}
