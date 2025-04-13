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
	if cfg.TelegramBotKey.IsNone() {
		return nil, nil
	}
	botToken := *cfg.TelegramBotKey.Unwrap()
	receivers := *cfg.TelegramReceiverIDs.UnwrapOr([]int64{})
	if cfg.TelegramReceiverIDs.IsNone() && len(receivers) == 0 {
		return nil, errors.New("telegram bot key was set but no receiver is set")
	}
	handler, err := telegram.New(botToken)
	if err != nil {
		return nil, err
	}
	handler.AddReceivers(receivers...)

	return handler, nil
}
