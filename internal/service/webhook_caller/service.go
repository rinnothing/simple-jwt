package webhook

import (
	"encoding/json"
	"fmt"

	"github.com/rinnothing/simple-jwt/internal/config"

	"go.uber.org/zap"
	"resty.dev/v3"
)

type WebhookService interface {
	CallWebhook(ip string) error
}

type WebhookServiceImpl struct {
	l *zap.Logger

	client *resty.Client
}

func NewService(cfg config.WebhookConfig, l *zap.Logger) WebhookService {
	return &WebhookServiceImpl{
		l:      l,
		client: resty.New().SetBaseURL(cfg.HttpAddress).SetRetryCount(cfg.RetryCount),
	}
}

type PostRequest struct {
	Message string `json:"message"`
}

func (w *WebhookServiceImpl) CallWebhook(ip string) error {
	request := w.client.R()

	requestBytes, err := json.Marshal(PostRequest{Message: fmt.Sprintf("try to auth from unknow ip found: %s", ip)})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := request.
		SetHeader("Content-Type", "application/json").
		SetBody(requestBytes).
		Post("/")

	if err != nil {
		return fmt.Errorf("failed to call webhook: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("request to webhook failed: HTTP %d  body: %s", resp.StatusCode(), resp.String())
	}

	return nil
}
