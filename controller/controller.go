package controller

import (
	"context"

	"go.uber.org/zap"

	"github.com/fmotalleb/the-one/config"
	"github.com/fmotalleb/the-one/logging"
	"github.com/fmotalleb/the-one/process"
)

var log = logging.LazyLogger("controller")

func Boot(ctx context.Context, cfg *config.Config) error {
	l := log().Named("Boot")
	l.Info("booting service controller")

	rootServices, err := cfg.BuildServiceGraph()
	if err != nil {
		return err
	}

	// Compile Templates
	if err = compileTemplates(cfg, l); err != nil {
		return err
	}

	for _, svc := range rootServices {
		svc.TraverseNode(func(s config.ServiceNode) {
			s.Data.OnDependChange = func() {
				if s.Data.GetDependCount() == 1 {
					executeSvcNode(ctx, s)
				} else {
					l.Debug("server requirements unmet", zap.Any("service", s.Data.Name()))
				}
			}
		})
	}
	for _, svc := range rootServices {
		executeSvcNode(ctx, svc)
	}
	select {}
	return nil
}

func executeSvcNode(ctx context.Context, svc config.ServiceNode) {
	s := svc.Data
	proc := process.New(s)
	go proc.Execute(ctx)
	go func() {
		proc.WaitForHealthy()
		svc.Traverse(config.ReduceDependCount)
	}()
}
