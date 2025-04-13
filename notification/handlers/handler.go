package handlers

import (
	"github.com/nikoksr/notify"

	"github.com/fmotalleb/the-one/config"
)

type NotifierBuilder = func(config.ContactPoint) (notify.Notifier, error)

var constructors = make([]NotifierBuilder, 0)

func FindHandler(cfg config.ContactPoint) ([]notify.Notifier, error) {
	results := make([]notify.Notifier, 0)
	for _, builder := range constructors {
		notifier, err := builder(cfg)
		if err != nil {
			return nil, err
		} else if notifier != nil {
			results = append(results, notifier)
		}
	}
	return results, nil
}
