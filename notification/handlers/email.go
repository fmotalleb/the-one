package handlers

import (
	"cmp"
	"errors"
	"strings"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/mail"

	"github.com/fmotalleb/the-one/config"
)

func init() {
	constructors = append(constructors, emailHandler)
}

func emailHandler(cfg config.ContactPoint) (notify.Notifier, error) {
	if cfg.SMTPHost == nil {
		return nil, nil
	}
	smtpHost := *cfg.SMTPHost

	smtpHostName := cmp.Or(*cfg.SMTPHostName, strings.Split(smtpHost, ":")[0])
	smtpUser := *cfg.SMTPUser.Unwrap()
	smtpPass := *cfg.SMTPPass.Unwrap()
	smtpFrom := cfg.SMTPFrom.UnwrapOr(smtpUser)
	smtpReceivers := make([]string, len(cfg.SMTPReceivers))
	for index, i := range cfg.SMTPReceivers {
		smtpReceivers[index] = *i.Unwrap()
	}

	if len(smtpReceivers) == 0 {
		return nil, errors.New("email host was set but missing smtp_receivers")
	}

	handler := mail.New(smtpFrom, smtpHost)
	// url, err := url.Parse(smtpHost)
	// if err != nil {
	// 	return nil, err
	// }

	handler.AuthenticateSMTP("", smtpUser, smtpPass, smtpHostName)

	handler.AddReceivers(smtpReceivers...)

	return handler, nil
}
