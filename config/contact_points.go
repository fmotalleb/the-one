package config

import (
	"net/http"

	"github.com/fmotalleb/the-one/types/option"
)

type ContactPoint struct {
	Name                option.Some[string]       `mapstructure:"name,omitempty"  yaml:"name"`
	TelegramBotKey      option.OptionalT[string]  `mapstructure:"telegram,omitempty"  yaml:"telegram"`
	TelegramReceiverIDs []option.OptionalT[int64] `mapstructure:"telegram_receivers,omitempty"  yaml:"telegram_receivers"`

	HTTPWebhookAddress option.OptionalT[string] `mapstructure:"webhook,omitempty"  yaml:"webhook"`

	// Generic http flags
	HTTPMethod      option.OptionalT[string]              `mapstructure:"http_method,omitempty"  yaml:"http_method"`
	HTTPContentType option.OptionalT[string]              `mapstructure:"http_content_type,omitempty"  yaml:"http_content_type"`
	HTTPHeaders     map[string][]option.OptionalT[string] `mapstructure:"http_headers,omitempty"  yaml:"http_headers"`

	// Mail Configs
	SMTPHost      option.OptionalT[string]   `mapstructure:"smtp_host,omitempty"  yaml:"smtp_host"`
	SMTPHostName  option.OptionalT[string]   `mapstructure:"smtp_hostname,omitempty"  yaml:"smtp_hostname"`
	SMTPUser      option.OptionalT[string]   `mapstructure:"smtp_user,omitempty"  yaml:"smtp_user"`
	SMTPPass      option.OptionalT[string]   `mapstructure:"smtp_pass,omitempty"  yaml:"smtp_pass"`
	SMTPFrom      option.OptionalT[string]   `mapstructure:"smtp_from,omitempty"  yaml:"smtp_from"`
	SMTPReceivers []option.OptionalT[string] `mapstructure:"smtp_receivers,omitempty"  yaml:"smtp_receivers"`
}

func (c ContactPoint) GetHTTPHeaders() http.Header {
	result := http.Header{}
	if len(c.HTTPHeaders) == 0 {
		return result
	}
	for key, values := range c.HTTPHeaders {
		realValues := make([]string, len(values))
		for index, item := range values {
			realValues[index] = *item.Unwrap()
		}
		result[key] = realValues
	}
	return result
}
