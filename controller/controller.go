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
					proc := process.New(s)
					go proc.Execute(ctx)
					proc.WaitForHealthy()
					svc.Traverse(config.ReduceDependCount)
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
