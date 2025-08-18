package config

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"github.com/fmotalleb/go-tools/writer"

	"github.com/fmotalleb/the-one/types/option"
)

// Service represents a single service definition in the system,
// including metadata, execution details, lifecycle, and dependencies.
type Service struct {
	// NameValue is the unique name of the service.
	// This field is required.
	NameValue string `mapstructure:"name,omitempty" yaml:"name" validate:"required"`

	// Enabled specifies whether the service is in the service tree or not.
	// If false, the service will be ignored.
	// TODO: switch with custom default true bool replacement
	Enabled option.OptionalT[bool] `mapstructure:"enabled,omitempty" yaml:"enabled"`

	// An absolute path to executable binary.
	// This field is required.
	Executable string `mapstructure:"executable,omitempty" yaml:"executable" validate:"required"`

	// Arguments is the list of optional arguments passed to the executable.
	Arguments []string `mapstructure:"args,omitempty" yaml:"args"`

	// Environments is a map of environment variables passed to the process.
	// Values can be explicitly unset or null, if inherit is set to false or unset process will be started with zero environment variable.
	Environments map[string]string `mapstructure:"env,omitempty" yaml:"env"`

	// Acts like list of .env files for the service
	EnvironmentFile []string `mapstructure:"env_file,omitempty" yaml:"env_file"`

	// Passes environment variables of the init system to child service.
	// Unset or null acts like false.
	InheritEnviron bool `mapstructure:"inherit_env,omitempty" yaml:"inherit_env"`

	// WorkingDir defines the working directory for the service process.
	// If unset, it defaults to the current working directory of the init system.
	WorkingDir string `mapstructure:"working_dir,omitempty" yaml:"working_dir" validate:"workingdir"`

	// ProcessCount specifies how many instances of the executable to run.
	// Defaults to 1 if not set.
	ProcessCount int `mapstructure:"process_count,omitempty" yaml:"process_count"`

	// Restart holds the configuration for automatic restarts on failure.
	// If unset, will use default restart behavior:
	// - Min Delay: 1s
	// - Max Delay: 15s
	// - Count: None
	Restart RestartConfig `mapstructure:"restart,omitempty" yaml:"restart"`

	// Timeout is the maximum time allowed for starting or stopping the process.
	// A zero or unset value means no timeout is enforced.
	// Its Considered in *-shot based services, in normal services this field has no means.
	Timeout time.Duration `mapstructure:"timeout,omitempty" yaml:"timeout"`

	// Type determines the kind of service [ServiceType]
	Type ServiceType `mapstructure:"type,omitempty" yaml:"type"`

	// Lazy indicates whether the service should be started lazily,
	// i.e., only when required by a dependent.
	Lazy bool `mapstructure:"lazy,omitempty" yaml:"lazy"`

	// HealthCheck defines the periodic check configuration to validate service health.
	HealthCheck HealthCheckConfig `mapstructure:"health_check,omitempty" yaml:"health_check"`

	// Requirements is a list of service names that must be successfully started before this one.
	Requirements []string `mapstructure:"requires,omitempty" yaml:"requires" validate:"dive,required"`

	//! **Dropped due to being ambiguous.**
	// DependencyItems lists services that must be stopped before this one starts.
	// These are soft constraints used in sequencing, not hard dependencies.
	// Dependencies []option.Optional[string] `mapstructure:"after,omitempty" yaml:"after"`

	// RequiredBy are services that depend on this one.
	RequiredBy []string `mapstructure:"dependents,omitempty" yaml:"dependents" validate:"dive,required"`

	// TODO: Still in process of freezing the configuration
	// Currently needs a slice
	// [type,parameter]
	// [stdout]
	// [stderr]
	// [file,./test.log]
	StdOut *writer.Writer `mapstructure:"stdout,omitempty" yaml:"stdout"`

	// By default will use [StdOut] if not provided
	StdErr *writer.Writer `mapstructure:"stderr,omitempty" yaml:"stderr"`

	dependCount atomic.Int64

	//! TODO its sloppy and temporary (final result is like this anyway)
	OnDependChange func()
}

func (s *Service) Name() string {
	return s.NameValue
}

func (s *Service) Dependencies() []string {
	return s.Requirements
}

func (s *Service) Dependents() []string {
	return s.RequiredBy
}

func (s *Service) GetType() ServiceType {
	return OngoingService
}

func (s *Service) GetProcessCount() int {
	return DefaultProcessCount
}

func (s *Service) GetRestart() RestartConfig {
	return DefaultRestartConfig
}

func (s *Service) GetOut() io.Writer {
	if s.StdOut != nil {
		return s.StdOut
	}
	return writer.NewStdErr()
}

func (s *Service) GetDependCount() int64 {
	return s.dependCount.Load()
}

func (s *Service) GetErr() io.Writer {
	if s.StdErr != nil {
		return s.StdErr
	}
	return s.GetOut()
}

func (s *Service) String() string {
	return fmt.Sprintf("%s %d", s.Name(), s.GetDependCount())
}

func IncreaseDependCount(s *Service) {
	s.dependCount.Add(1)
}

func ReduceDependCount(s *Service) {
	if s.GetDependCount() > 1 {
		s.dependCount.Add(-1)
		s.OnDependChange()
	}
}
