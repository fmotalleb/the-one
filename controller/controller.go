package controller

import (
	"context"
	"errors"
	"os"

	"go.uber.org/zap"

	"github.com/fmotalleb/the-one/config"
	"github.com/fmotalleb/the-one/logging"
	"github.com/fmotalleb/the-one/process"
	"github.com/fmotalleb/the-one/renderer"
)

var log = logging.LazyLogger("controller")

func Boot(ctx context.Context, cfg *config.Config) error {
	l := log().Named("Boot")
	l.Info("booting service controller")
	// Compile Templates
	for _, t := range cfg.Templates {
		tl := l.With(zap.String("src", t.GetSourceDirectory()))
		tl.Debug(
			"rendering template directory",
		)
		if err := renderer.RenderTemplates(&t); err != nil && t.GetIsFatal() {
			tl.Error(
				"failed to render templates",
				zap.Error(err),
			)
			return errors.Join(
				ErrEngineBoot,
				err,
			)
		}
	}

	if len(cfg.Templates) != 0 {
		l.Debug("finished template rendering")
	} else {
		l.Debug("no template was found, ignoring the step")
	}
	for _, svc := range cfg.Services {
		// Currently does nothing but will be used to
		ctrl := make(chan process.ServiceMessage)
		signal := make(chan process.ServiceMessage)
		go trackChannel(signal)
		std := process.StdConfig{
			In:  os.Stdin,
			Out: svc.GetOut(),
			Err: svc.GetErr(),
		}
		mgr := process.NewServiceManager(ctx, &svc, ctrl, signal, std)
		err := mgr.Start()
		l.Error("failed to start process", zap.Error(err))
	}
	select {}
	return nil
}

func trackChannel(ch chan process.ServiceMessage) {
	l := log().Named("Tracker")
	for v := range ch {
		l.Debug("message received", zap.Any("signal", v))
	}
}
