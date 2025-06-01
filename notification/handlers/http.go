package handlers

import (
	stdhttp "net/http"
	"time"

	"github.com/nikoksr/notify"

	"github.com/fmotalleb/the-one/config"

	"github.com/nikoksr/notify/service/http"
)

func init() {
	constructors = append(constructors, httpWebhookHandler)
}

func httpWebhookHandler(cfg config.ContactPoint) (notify.Notifier, error) {
	if cfg.HTTPWebhookAddress.IsSome() {
		wh := *cfg.HTTPWebhookAddress.Unwrap()
		client := http.New()
		method := cfg.HTTPMethod.UnwrapOr(stdhttp.MethodPost)
		client.AddReceivers(&http.Webhook{
			URL:          wh,
			Header:       cfg.GetHTTPHeaders(),
			Method:       method,
			ContentType:  "application/json",
			BuildPayload: generatePayload,
		})
		return client, nil
	}
	return nil, nil
}

func generatePayload(subject, message string) any {
	payload := map[string]any{
		"subject": subject,
		"message": message,
		"time":    time.Now().UTC().Format(time.RFC3339),
	}
	return payload
}
