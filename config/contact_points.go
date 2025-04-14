package config

import (
	"net/http"

	"github.com/fmotalleb/the-one/types/option"
)

type ContactPoint struct {
	Name                option.Some[string]       `mapstructure:"name,omitempty"  yaml:"name"`
	TelegramBotKey      option.OptionalT[string]  `mapstructure:"telegram,omitempty"  yaml:"telegram"`
	TelegramReceiverIDs option.OptionalT[[]int64] `mapstructure:"telegram_receivers,omitempty"  yaml:"telegram_receivers"`

	HTTPWebhookAddress option.OptionalT[string] `mapstructure:"webhook,omitempty"  yaml:"webhook"`

	// Generic http flags
	HTTPMethod      option.OptionalT[string]              `mapstructure:"http_method,omitempty"  yaml:"http_method"`
	HTTPContentType option.OptionalT[string]              `mapstructure:"http_content_type,omitempty"  yaml:"http_content_type"`
	HTTPHeaders     option.OptionalT[map[string][]string] `mapstructure:"http_headers,omitempty"  yaml:"http_headers"`
}

func (c ContactPoint) GetHTTPHeaders() http.Header {
	if c.HTTPHeaders.IsNone() {
		return http.Header{}
	}
	headers := *c.HTTPHeaders.Unwrap()
	return headers
}
