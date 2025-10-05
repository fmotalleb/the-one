package handlers

import (
	"errors"
	"strings"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/mail"

	"github.com/fmotalleb/go-tools/ptrcmp"

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

	smtpHostName := ptrcmp.Or(cfg.SMTPHostName, strings.Split(smtpHost, ":")[0])
	if cfg.SMTPUser == nil || cfg.SMTPPass == nil || cfg.SMTPReceivers == nil || len(*cfg.SMTPReceivers) == 0 {
		return nil, errors.New("missing required SMTP parameters")
	}
	smtpUser := *cfg.SMTPUser
	smtpPass := *cfg.SMTPPass
	smtpFrom := ptrcmp.Or(cfg.SMTPFrom, smtpUser)
	smtpReceivers := *cfg.SMTPReceivers

	handler := mail.New(smtpFrom, smtpHost)

	handler.AuthenticateSMTP("", smtpUser, smtpPass, smtpHostName)

	handler.AddReceivers(smtpReceivers...)

	return handler, nil
}
