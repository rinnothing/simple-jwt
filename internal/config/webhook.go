package config

type WebhookConfig struct {
	HttpAddress string `yaml:"http_address" env:"WEBHOOK_HTTP_ADDRESS"`
	RetryCount  int    `yaml:"retry_count" env:"WEBHOOK_RETRY_COUNT"`
}
