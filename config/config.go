package config

import (
	"github.com/fmotalleb/the-one/logging"
	"github.com/fmotalleb/the-one/types/option"
)

var log = logging.LazyLogger("core.config")

type Config struct {
	Services  []Service      `mapstructure:"services,omitempty"`
	Templates []Template     `mapstructure:"templates,omitempty"`
	Contacts  []ContactPoint `mapstructure:"contacts,omitempty"`
}

func (c *Config) GetServices() ([]Service, error) {
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
		after := option.WrapAll(dependencies)
		s.After = after
		s.Dependents = []option.Optional[string]{}
		services[index] = s
	}
	return services
}

func generateDependencyList(c *Config) map[string][]string {
	dependencies := map[string][]string{}
	for _, s := range c.Services {
		name := *s.Name.Unwrap()
		after := option.UnwrapAll(s.After)
		dependencies[name] = append(dependencies[name], after...)

		dependents := option.UnwrapAll(s.Dependents)
		for _, d := range dependents {
			dependencies[d] = append(dependencies[d], name)
		}
	}
	return dependencies
}
