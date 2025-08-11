package controller

import (
	"context"

	"github.com/fmotalleb/go-tools/tree"

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
	weightServices(rootServices)

	// Compile Templates
	if err = compileTemplates(cfg, l); err != nil {
		return err
	}

	for _, svc := range rootServices {
		svc.TraverseNode(func(s *tree.Node[*config.Service]) {
			s.Data.OnDependChange = func() {
				if s.Data.GetDependCount() == 1 {
					executeSvcNode(ctx, s)
				} else {
					println(s.Data.String())
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

func executeSvcNode(ctx context.Context, svc *tree.Node[*config.Service]) {
	s := svc.Data
	proc := process.New(s)
	go proc.Execute(ctx)
	go func() {
		proc.WaitForHealthy()
		svc.Traverse(config.ReduceDependCount)
	}()
}

func weightServices(rootServices []*tree.Node[*config.Service]) {
	for _, root := range rootServices {
		weightServices(root.Children())
		if root.Data.GetDependCount() == 0 {
			root.Traverse(config.IncreaseDependCount)
		}
	}
}
