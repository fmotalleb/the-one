package config

import (
	"time"

	"github.com/fmotalleb/the-one/types/option"
)

type RestartConfig struct {
	Count       option.Optional[uint]          `mapstructure:"count,omitempty"`
	Delay       option.Optional[time.Duration] `mapstructure:"delay,omitempty"`
	DelayMax    option.Optional[time.Duration] `mapstructure:"delay_max,omitempty"`
	OkExitCodes option.Optional[[]int]         `mapstructure:"ok_exit_codes,omitempty"`
}

// GetCount returns an unsigned-integer value and a boolean.
// if the count is omitted the boolean is false and indicates
// the service should not be stopped if failed multiple times.
func (r *RestartConfig) GetCount() (uint, bool) {
	if r.Count.IsSome() {
		return *r.Count.Unwrap(), true
	}
	return 0, false
}

// GetDelayBegin returns value of minimum allowed delay set in config.
// if omitted it will return [DefaultRestartDelayBegin].
func (r *RestartConfig) GetDelayBegin() time.Duration {
	return *r.Delay.UnwrapOr(DefaultRestartDelayBegin)
}

// GetDelayMax returns value of maximum allowed delay set in config.
// if omitted it will return [DefaultRestartDelayMax].
func (r *RestartConfig) GetDelayMax() time.Duration {
	return *r.DelayMax.UnwrapOr(DefaultRestartDelayMax)
}

// GetExitCodes returns value of ok exit codes set in config.
// if omitted it will return [DefaultRestartOkCodes] [0].
func (r *RestartConfig) GetExitCodes() []int {
	return *r.OkExitCodes.UnwrapOr(DefaultRestartOkCodes)
}
