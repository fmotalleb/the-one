package config

import (
	"io"
	"time"

	"github.com/fmotalleb/go-tools/writer"

	"github.com/fmotalleb/the-one/types/option"
)

// Service represents a single service definition in the system,
// including metadata, execution details, lifecycle, and dependencies.
type Service struct {
	// Name is the unique name of the service.
	// This field is required.
	Name option.Some[string] `mapstructure:"name,omitempty" yaml:"name"`

	// Enabled specifies whether the service is in the service tree or not.
	// If false, the service will be ignored.
	Enabled option.OptionalT[bool] `mapstructure:"enabled,omitempty" yaml:"enabled"`

	// An absolute path to executable binary.
	// This field is required.
	Executable option.Some[string] `mapstructure:"executable,omitempty" yaml:"executable"`

	// Arguments is the list of optional arguments passed to the executable.
	Arguments []option.OptionalT[string] `mapstructure:"args,omitempty" yaml:"args"`

	// Environments is a map of environment variables passed to the process.
	// Values can be explicitly unset or null, if inherit is set to false or unset process will be started with zero environment variable.
	Environments map[string]option.Optional[string] `mapstructure:"env,omitempty" yaml:"env"`

	// Acts like list of .env files for the service
	EnvironmentFile []option.Optional[string] `mapstructure:"env_file,omitempty" yaml:"env_file"`

	// Passes environment variables of the init system to child service.
	// Unset or null acts like false.
	InheritEnviron option.Optional[bool] `mapstructure:"inherit_env,omitempty" yaml:"inherit_env"`

	// WorkingDir defines the working directory for the service process.
	// If unset, it defaults to the current working directory of the init system.
	WorkingDir option.OptionalT[string] `mapstructure:"working_dir,omitempty" yaml:"working_dir"`

	// ProcessCount specifies how many instances of the executable to run.
	// Defaults to 1 if not set.
	ProcessCount option.OptionalT[int] `mapstructure:"process_count,omitempty" yaml:"process_count"`

	// Restart holds the configuration for automatic restarts on failure.
	// If unset, will use default restart behavior:
	// - Min Delay: 1s
	// - Max Delay: 15s
	// - Count: None
	Restart option.Optional[RestartConfig] `mapstructure:"restart,omitempty" yaml:"restart"`

	// Timeout is the maximum time allowed for starting or stopping the process.
	// A zero or unset value means no timeout is enforced.
	// Its Considered in *-shot based services, in normal services this field has no means.
	Timeout option.OptionalT[time.Duration] `mapstructure:"timeout,omitempty" yaml:"timeout"`

	// Type determines the kind of service [ServiceType]
	Type option.Optional[ServiceType] `mapstructure:"type,omitempty" yaml:"type"`

	// Lazy indicates whether the service should be started lazily,
	// i.e., only when required by a dependent.
	Lazy option.Optional[bool] `mapstructure:"lazy,omitempty" yaml:"lazy"`

	// HealthCheck defines the periodic check configuration to validate service health.
	HealthCheck option.Optional[HealthCheckConfig] `mapstructure:"health_check,omitempty" yaml:"health_check"`

	// Requirements is a list of service names that must be successfully started before this one.
	Requirements []option.Optional[string] `mapstructure:"requires,omitempty" yaml:"requires"`

	// After lists services that must be stopped before this one starts.
	// These are soft constraints used in sequencing, not hard dependencies.
	After []option.Optional[string] `mapstructure:"After,omitempty" yaml:"After"`

	// Dependents are services that depend on this one.
	// Internally, this is translated to `After` entries in those dependent services.
	// This field is cleared before execution.
	Dependents []option.Optional[string] `mapstructure:"dependents,omitempty" yaml:"dependents"`

	StdOut *writer.Writer `mapstructure:"stdout,omitempty" yaml:"stdout"`
	StdErr *writer.Writer `mapstructure:"stderr,omitempty" yaml:"stderr"`
}

func (s *Service) GetName() string {
	return s.Name.UnwrapOr("")
}

func (s *Service) GetType() ServiceType {
	return s.Type.UnwrapOr(OngoingService)
}

func (s *Service) GetProcessCount() int {
	return s.ProcessCount.UnwrapOr(DefaultProcessCount)
}

func (s *Service) GetRestart() RestartConfig {
	return s.Restart.UnwrapOr(DefaultRestartConfig)
}

func (s *Service) GetOut() io.Writer {
	if s.StdOut != nil {
		return s.StdOut
	}
	return writer.NewStdErr()
}

func (s *Service) GetErr() io.Writer {
	if s.StdErr != nil {
		return s.StdErr
	}
	return s.GetOut()
}
