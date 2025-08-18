package controller

import (
	"context"

	"go.uber.org/zap"

	"github.com/fmotalleb/go-tools/log"
	"github.com/fmotalleb/go-tools/tree"

	"github.com/fmotalleb/the-one/config"
	"github.com/fmotalleb/the-one/process"
)

func Boot(ctx context.Context, cfg *config.Config) error {
	l := log.Of(ctx)
	l = l.Named("Boot")
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
		svc.TraverseNode(func(s *tree.Node[*config.Service]) {
			s.Data.OnDependChange = func() {
				depCount := s.Data.GetDependCount()
				if depCount == 1 {
					executeSvcNode(ctx, s)
				} else if depCount > 1 {
					l.Debug("service requirements unmet", zap.Any("service", s.Data.Name()))
				}
			}
		})
	}
	for _, svc := range rootServices {
		l.Debug("starting srv", zap.Any("srv", svc.Data))
		executeSvcNode(ctx, svc)
	}
	<-make(chan int)
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
