package config

import (
	"time"

	"github.com/fmotalleb/the-one/types/option"
)

type HealthCheckConfig struct {
	Type     string                          `mapstructure:"type,omitempty" yaml:"type"`       // "http", "tcp", "cmd"
	Address  option.OptionalT[string]        `mapstructure:"address,omitempty" yaml:"address"` // for http/tcp
	Command  option.Optional[[]string]       `mapstructure:"command,omitempty" yaml:"command"` // for cmd
	Interval option.OptionalT[time.Duration] `mapstructure:"interval,omitempty" yaml:"interval"`
	Timeout  option.OptionalT[time.Duration] `mapstructure:"timeout,omitempty" yaml:"timeout"`
	Retries  option.OptionalT[int]           `mapstructure:"retries,omitempty" yaml:"retries"`
}
