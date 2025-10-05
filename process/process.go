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

	"github.com/sethvargo/go-retry"
)

type Process struct {
	name  string
	exe   string
	args  []string
	retry retry.Backoff
}

func New(svc *config.Service) *Process {
	retryCfg := svc.Retry

	r := retry.NewFibonacci(retryCfg.GetDelayBegin())
	r = retry.WithMaxDuration(retryCfg.GetDelayMax(), r)
	if count, ok := retryCfg.GetCount(); ok {
		r = retry.WithMaxRetries(uint64(count), r)
	}
	return &Process{
		name:  svc.Name(),
		exe:   svc.Executable,
		args:  svc.Arguments,
		retry: r,
	}
}

func (p *Process) Execute(ctx context.Context) error {
	l := log.Of(ctx).Named("process.Execute")

	l.Warn(p.name, zap.String("state", "online"))
	err := retry.Do(ctx, p.retry, p.spawnProcess)
	return err
}

func (p *Process) spawnProcess(ctx context.Context) error {
	return nil
}

func (p *Process) WaitForHealthy() error {
	time.Sleep(time.Second)
	return nil
}
