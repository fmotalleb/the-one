package config

import (
	"net/http"

	"github.com/fmotalleb/the-one/types/option"
)

type ContactPoint struct {
	TelegramBotKey      option.Optional[string] `mapstructure:"telegram,omitempty"`
	TelegramReceiverIDs option.Option[[]int64]  `mapstructure:"telegram_receivers,omitempty"`

	HTTPWebhookAddress option.Optional[string] `mapstructure:"webhook,omitempty"`

	// Generic http flags
	HTTPMethod      option.Optional[string]              `mapstructure:"http_method,omitempty"`
	HTTPContentType option.Optional[string]              `mapstructure:"http_method,omitempty"`
	HTTPHeaders     option.Optional[map[string][]string] `mapstructure:"http_headers,omitempty"`
}

func (c ContactPoint) GetHTTPHeaders() http.Header {
	if c.HTTPHeaders.IsNone() {
		return http.Header{}
	}
	headers := *c.HTTPHeaders.Unwrap()
	return headers
}
