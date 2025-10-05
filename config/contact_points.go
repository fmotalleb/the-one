package config

import (
	"net/http"
)

type ContactPoint struct {
	Name *string `mapstructure:"name,omitempty" yaml:"name" validate:"required"`

	TelegramBotKey      *string  `mapstructure:"telegram,omitempty"  yaml:"telegram"`
	TelegramReceiverIDs *[]int64 `mapstructure:"telegram_receivers,omitempty"  yaml:"telegram_receivers"`

	HTTPWebhookAddress *string `mapstructure:"webhook,omitempty"  yaml:"webhook"`

	// Generic http flags
	HTTPMethod      *string              `mapstructure:"http_method,omitempty"  yaml:"http_method"`
	HTTPContentType *string              `mapstructure:"http_content_type,omitempty"  yaml:"http_content_type"`
	HTTPHeaders     *map[string][]string `mapstructure:"http_headers,omitempty"  yaml:"http_headers"`

	// Mail Configs
	SMTPHost      *string   `mapstructure:"smtp_host,omitempty"  yaml:"smtp_host"`
	SMTPHostName  *string   `mapstructure:"smtp_hostname,omitempty"  yaml:"smtp_hostname"`
	SMTPUser      *string   `mapstructure:"smtp_user,omitempty"  yaml:"smtp_user"`
	SMTPPass      *string   `mapstructure:"smtp_pass,omitempty"  yaml:"smtp_pass"`
	SMTPFrom      *string   `mapstructure:"smtp_from,omitempty"  yaml:"smtp_from"`
	SMTPReceivers *[]string `mapstructure:"smtp_receivers,omitempty"  yaml:"smtp_receivers"`
}

func (c ContactPoint) GetHTTPHeaders() http.Header {
	result := http.Header{}
	if c.HTTPHeaders == nil {
		return result
	}
	for key, values := range *c.HTTPHeaders {
		realValues := make([]string, len(values))
		copy(realValues, values)
		result[key] = realValues
	}
	return result
}
