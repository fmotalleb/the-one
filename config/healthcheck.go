package config

import (
	"time"
)

type HealthCheckConfig struct {
	Type          string         `mapstructure:"type,omitempty" yaml:"type"`       // "http", "tcp", "cmd"
	Address       *string        `mapstructure:"address,omitempty" yaml:"address"` // for http/tcp
	Command       *[]string      `mapstructure:"command,omitempty" yaml:"command"` // for cmd
	ResultMatcher *string        `mapstructure:"output_matcher,omitempty" yaml:"output_matcher"`
	Interval      *time.Duration `mapstructure:"interval,omitempty" yaml:"interval"`
	Timeout       *time.Duration `mapstructure:"timeout,omitempty" yaml:"timeout"`
	Retries       *int           `mapstructure:"retries,omitempty" yaml:"retries"`
	OkExitCodes   *[]int         `mapstructure:"ok_exit_codes,omitempty" yaml:"ok_exit_codes"`
}
