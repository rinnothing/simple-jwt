package config

type WebhookConfig struct {
	HttpAddress string `json:"http_address" env:"WEBHOOK_HTTP_ADDRESS"`
	RetryCount  int    `json:"retry_count" env:"WEBHOOK_RETRY_COUNT"`
}
