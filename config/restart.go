package config

import (
	"math"
	"time"

	"github.com/fmotalleb/the-one/types/option"
)

type RestartConfig struct {
	// Count is the number of times to restart the service if it fails.
	// If omitted, the service will be restarted indefinitely.
	// If set to 0, the service will not be restarted.
	Count option.Optional[uint] `mapstructure:"count,omitempty" yaml:"count"`

	// Delay is the minimum delay between restarts.
	// If omitted, the default is 1 second.
	// If set to 0, the service will be restarted immediately.
	Delay option.Optional[time.Duration] `mapstructure:"delay,omitempty" yaml:"delay"`

	// DelayMax is the maximum delay between restarts.
	// If omitted, the default is 16 seconds.
	// If set to 0, the service will be restarted immediately.
	// If set to a value less than Delay, considered as Delay.
	// Each restart will increase the delay by a factor of 2, up to DelayMax.
	// For example, if Delay is 1 second and DelayMax is 16 seconds,
	// the delays will be 1, 2, 4, 8, and then 16 seconds and will stay 16seconds until.
	DelayMax option.Optional[time.Duration] `mapstructure:"delay_max,omitempty" yaml:"delay_max"`
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

func (r *RestartConfig) GetDelay(iteration uint) (time.Duration, bool) {
	count, ok := r.GetCount()
	if ok && count <= iteration || iteration > restartAbsoluteMax {
		return 0, false
	}

	maxDelay := r.GetDelayMax()

	if iteration >= restartMaxCalculableIteration {
		return maxDelay, true
	}
	multiplier := math.Pow(restartDelayPowerBase, float64(iteration))
	delay := r.GetDelayBegin() * time.Duration(multiplier)

	if delay > maxDelay {
		return maxDelay, true
	}
	return delay, true
}
