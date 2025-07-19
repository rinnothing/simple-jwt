package webhook

type WebhookService interface {
	CallWebhook(ip string) error
}
