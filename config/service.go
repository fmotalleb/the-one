package config

import (
	"time"

	"github.com/fmotalleb/the-one/types/option"
)

type Service struct {
	// service metadata
	Name option.Some[string] `mapstructure:"name,omitempty" yaml:"name"`
	// Description option.Optional[string] `mapstructure:"description,omitempty"`

	Enabled option.OptionalT[bool] `mapstructure:"enabled,omitempty" yaml:"enabled"`

	// process information
	Executable   option.Some[string]                `mapstructure:"executable,omitempty" yaml:"executable"`
	Arguments    option.Optional[[]string]          `mapstructure:"args,omitempty" yaml:"args"`
	Environments option.Optional[map[string]string] `mapstructure:"env,omitempty" yaml:"env"`
	WorkingDir   option.OptionalT[[]string]         `mapstructure:"working_dir,omitempty" yaml:"working_dir"`
	ProcessCount option.OptionalT[int]              `mapstructure:"process_count,omitempty" yaml:"process_count"`

	// process management
	Restart     option.Optional[RestartConfig]  `mapstructure:"restart,omitempty" yaml:"restart"`
	Timeout     option.OptionalT[time.Duration] `mapstructure:"timeout,omitempty" yaml:"timeout"`
	Type        option.Optional[ServiceType]    `mapstructure:"type,omitempty" yaml:"type"`
	Lazy        option.Optional[bool]           `mapstructure:"lazy,omitempty" yaml:"lazy"`
	OkExitCodes option.OptionalT[[]int]         `mapstructure:"ok_exit_codes,omitempty" yaml:"ok_exit_codes"`

	// dependency management
	Requirements option.Optional[[]string] `mapstructure:"requires,omitempty" yaml:"requires"`
	After        option.Optional[[]string] `mapstructure:"After,omitempty" yaml:"After"`
	Dependents   option.Optional[[]string] `mapstructure:"dependents,omitempty" yaml:"dependents"`
}
