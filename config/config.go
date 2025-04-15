package config

import (
	"github.com/fmotalleb/the-one/logging"
	"github.com/fmotalleb/the-one/types/option"
)

var log = logging.LazyLogger("core.config")

type Config struct {
	Services []Service      `mapstructure:"services,omitempty"`
	Contacts []ContactPoint `mapstructure:"contacts,omitempty"`
}

func (c *Config) GetServices() ([]Service, error) {
	// l := log().Named("GetServices")
	// sm := map[string]Service{}

	return reshapeServices(c), nil
}

func reshapeServices(c *Config) []Service {
	dependencies := generateDependencyList(c)
	services := updateDependencyList(c, dependencies)
	return services
}

func updateDependencyList(c *Config, dependencies map[string][]string) []Service {
	services := make([]Service, len(c.Services))
	for index, s := range c.Services {
		dependencies := dependencies[*s.Name.Unwrap()]
		s.After = option.NewOptional(&dependencies)
		s.Dependents = option.NewOptional[[]string](nil)
		services[index] = s
	}
	return services
}

func generateDependencyList(c *Config) map[string][]string {
	dependencies := map[string][]string{}
	for _, s := range c.Services {
		name := *s.Name.Unwrap()

		dependencies[name] = append(dependencies[name], *s.After.UnwrapOr([]string{})...)

		for _, d := range *s.Dependents.UnwrapOr([]string{}) {
			dependencies[d] = append(dependencies[d], name)
		}
	}
	return dependencies
}
