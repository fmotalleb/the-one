package controller

import (
	"context"
	"regexp"

	"github.com/fmotalleb/go-tools/broadcast"
	"go.uber.org/zap"
)

type (
	EventSource  uint16
	EngineState  uint16
	ServiceState uint16
)

// Values of this constants are important.
const (
	ESEngine  = EventSource(1)
	ESService = EventSource(2)

	EngineUp       = EngineState(0b001)
	EngineShutdown = EngineState(0b010)

	ServiceUp      = ServiceState(0b00001)
	ServiceAny     = ServiceState(0b00000)
	ServiceReady   = ServiceState(0b00010)
	ServiceStarted = ServiceState(0b00100) | ServiceUp
	ServiceHealthy = ServiceState(0b01000) | ServiceUp
	ServiceDown    = ServiceState(0b10000)
)

type Event struct {
	Source EventSource

	// Engine specific fields
	EngineState EngineState

	// Service specific fields
	ServiceState ServiceState
	ServiceName  string
}

func WaitForService(
	ctx context.Context,
	cast broadcast.Subscription[Event],
	matcher *regexp.Regexp,
	state ...ServiceState,
) {
	log := log().Named("WaitForService")

	listenOn := ServiceAny
	if len(state) != 0 {
		listenOn = state[0]
	}
	log.Debug("begin listening", zap.Uint16("listen-on", uint16(listenOn)))

	broadcast.Subscribe(
		cast,
		func(c <-chan Event) {
			for {
				select {
				case e := <-c:
					if e.Source == ESService &&
						matcher.FindIndex([]byte(e.ServiceName)) != nil &&
						listenOn&e.ServiceState == listenOn {
						log.Info("matched service event",
							zap.String("service", e.ServiceName),
							zap.Uint16("state", uint16(e.ServiceState)),
						)
						return
					} else {
						log.Debug("skipped event",
							zap.String("service", e.ServiceName),
							zap.Uint16("state", uint16(e.ServiceState)),
							zap.Any("event", e),
						)
					}
				case <-ctx.Done():
					log.Warn("context canceled or deadline exceeded")
					return
				}
			}
		},
	)
	log.Debug("finished waiting", zap.Uint16("listen-on", uint16(listenOn)))
}

func WaitForEngine(
	ctx context.Context,
	cast broadcast.Subscription[Event],
) {
	log := log().Named("WaitForEngine")
	listenOn := EngineUp
	log.Debug("begin listening", zap.Uint16("listen-on", uint16(listenOn)))

	broadcast.Subscribe(
		cast,
		func(c <-chan Event) {
			for {
				select {
				case e := <-c:
					if e.Source == ESEngine && listenOn&e.EngineState == listenOn {
						log.Info("matched engine event",
							zap.Uint16("state", uint16(e.EngineState)),
						)
						return
					} else {
						log.Debug("skipped event",
							zap.Uint16("state", uint16(e.EngineState)),
						)
					}
				case <-ctx.Done():
					log.Warn("context canceled or deadline exceeded")
					return
				}
			}
		},
	)
	log.Debug("finished waiting", zap.Uint16("listen-on", uint16(listenOn)))
}
