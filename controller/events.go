package controller

import (
	"context"
	"regexp"

	"github.com/fmotalleb/the-one/broadcast"
)

type EventSource uint16

const (
	ESEngine  = EventSource(0)
	ESService = EventSource(1)
)

type Event struct {
	Source    EventSource
	Parameter string
}

func WaitFor(
	ctx context.Context,
	cast broadcast.Subscription[Event],
	src EventSource,
	matcher regexp.Regexp,
) {
	notifier := cast.Subscribe()
	for {
		select {
		case e := <-notifier:
			if e.Source == src && matcher.FindIndex([]byte(e.Parameter)) != nil {
				break
			}
		case <-ctx.Done():
			break
		}
	}
}
