package config

import (
	"net/http"
)

type ContactPoint struct {
	Name *string `mapstructure:"name,omitempty" yaml:"name" validate:"required"`

	TelegramBotKey      *string  `mapstructure:"telegram,omitempty"  yaml:"telegram" validate:"omitempty,min=1"`
	TelegramReceiverIDs *[]int64 `mapstructure:"telegram_receivers,omitempty"  yaml:"telegram_receivers" validate:"omitempty,dive,required"`

	HTTPWebhookAddress *string              `mapstructure:"webhook,omitempty"  yaml:"webhook" validate:"omitempty,url"`
	HTTPMethod         *string              `mapstructure:"http_method,omitempty" default:"POST"  yaml:"http_method" validate:"omitempty,oneof=GET POST PUT PATCH DELETE"`
	HTTPContentType    *string              `mapstructure:"http_content_type,omitempty"  yaml:"http_content_type"`
	HTTPHeaders        *map[string][]string `mapstructure:"http_headers,omitempty"  yaml:"http_headers" validate:"omitempty,dive,keys,required,endkeys,dive,required"`

	// Mail Configs
	SMTPHost      *string   `mapstructure:"smtp_host,omitempty" yaml:"smtp_host"`
	SMTPHostName  *string   `mapstructure:"smtp_hostname,omitempty" yaml:"smtp_hostname" validate:"omitempty,hostname"`
	SMTPUser      *string   `mapstructure:"smtp_user,omitempty" yaml:"smtp_user" validate:"omitempty,min=1"`
	SMTPPass      *string   `mapstructure:"smtp_pass,omitempty" yaml:"smtp_pass" validate:"omitempty,min=1"`
	SMTPFrom      *string   `mapstructure:"smtp_from,omitempty" yaml:"smtp_from" validate:"omitempty,email"`
	SMTPReceivers *[]string `mapstructure:"smtp_receivers,omitempty" yaml:"smtp_receivers" validate:"omitempty,dive,email"`
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
