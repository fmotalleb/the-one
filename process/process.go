// TODO code in this package was a POC and must be rewritten in order to work
// Its just a mimic of what it should be
// Package process encapsulates logic behind process management and policies for individual processes
package process

import (
	"context"
	"time"

	"github.com/fmotalleb/go-tools/log"
	"go.uber.org/zap"

	"github.com/fmotalleb/the-one/config"
	"github.com/fmotalleb/the-one/helpers"
)

type Process struct {
	name string
	exe  string
	args []string
}

func New(svc *config.Service) *Process {
	return &Process{
		name: svc.Name(),
		exe:  *svc.Executable.Unwrap(),
		args: helpers.OptToSlice(svc.Arguments),
	}
}

func (p *Process) Execute(ctx context.Context) error {
	l := log.Of(ctx)
	l.Warn(p.name, zap.String("state", "online"))
	select {}
	return nil
}

func (p *Process) WaitForHealthy() error {
	time.Sleep(time.Second)
	return nil
}
