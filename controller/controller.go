package controller

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/fmotalleb/go-tools/tree"

	"github.com/fmotalleb/the-one/config"
	"github.com/fmotalleb/the-one/logging"
	"github.com/fmotalleb/the-one/process"
	"github.com/fmotalleb/the-one/renderer"
)

var log = logging.LazyLogger("controller")

func Boot(ctx context.Context, cfg *config.Config) error {
	l := log().Named("Boot")
	l.Info("booting service controller")

	rootServices, err := cfg.BuildServiceGraph()
	if err != nil {
		return err
	}
	weightServices(rootServices)

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
	for _, svc := range rootServices {
		svc.Traverse(func(s *config.Service) {
			s.OnDependChange = func() {
				if s.GetDependCount() == 0 {
					// Launch the entire start-wait-notify sequence in a new goroutine.
					// This handler now returns instantly, unblocking the main control flow.
					// go func() {
					proc := process.New(s)

					// Execute starts the process.
					go proc.Execute(ctx)

					// This now blocks only *inside this goroutine*.
					// It does NOT block other services from starting.
					proc.WaitForHealthy()

					// Once this service is healthy, notify its children to reduce their counts.
					// This will in turn trigger their OnDependChange handlers.

					svc.Traverse(config.ReduceDependCount)
					// }()
				}
			}
		})
	}
	for _, svc := range rootServices {
		svc.Traverse(config.ReduceDependCount)
	}
	select {}
	return nil
}

func weightServices(rootServices []*tree.Node[*config.Service]) {
	for _, root := range rootServices {
		root.Traverse(config.IncreaseDependCount)
		weightServices(root.Children())
	}
}
