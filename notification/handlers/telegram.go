package handlers

import (
	"errors"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"

	"github.com/fmotalleb/the-one/config"
)

func init() {
	constructors = append(constructors, telegramHandler)
}

func telegramHandler(cfg config.ContactPoint) (notify.Notifier, error) {
	if cfg.TelegramBotKey == nil {
		return nil, nil
	}
	botToken := *cfg.TelegramBotKey
	if cfg.TelegramReceiverIDs == nil || len(*cfg.TelegramReceiverIDs) == 0 {
		return nil, errors.New("telegram bot key passed but receiver list is null or empty")
	}
	receivers := *cfg.TelegramReceiverIDs

	handler, err := telegram.New(botToken)
	if err != nil {
		return nil, err
	}
	handler.AddReceivers(receivers...)

	return handler, nil
}
