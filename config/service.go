package config

import (
	"time"

	"github.com/fmotalleb/the-one/types/option"
)

type Service struct {
	// service metadata
	Name option.Some[string] `mapstructure:"name,omitempty"`
	// Description option.Optional[string] `mapstructure:"description,omitempty"`

	// process information
	Executable   option.Some[string]                `mapstructure:"executable,omitempty"`
	Arguments    option.Optional[[]string]          `mapstructure:"args,omitempty"`
	Environments option.Optional[map[string]string] `mapstructure:"env,omitempty"`
	WorkingDir   option.Optional[[]string]          `mapstructure:"working_dir,omitempty"`
	ProcessCount option.Optional[int]               `mapstructure:"process_count,omitempty"`

	// process management
	Restart option.Optional[RestartConfig] `mapstructure:"restart,omitempty"`
	Timeout option.Optional[time.Duration] `mapstructure:"timeout,omitempty"`
	Type    option.Optional[string]        `mapstructure:"type,omitempty"`
	Lazy    option.Optional[bool]          `mapstructure:"lazy,omitempty"`

	// dependency management
	Dependencies option.Optional[[]string] `mapstructure:"dependencies,omitempty"`
	Dependents   option.Optional[[]string] `mapstructure:"dependents,omitempty"`
}
