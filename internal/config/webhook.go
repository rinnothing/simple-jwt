package config

type WebhookConfig struct {
	HttpAddress string `yaml:"http_address"`
	RetryCount  int    `yaml:"retry_count"`
}
