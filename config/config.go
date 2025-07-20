package config

import "github.com/fmotalleb/go-tools/tree"

type Config struct {
	Services  []*Service     `mapstructure:"services,omitempty"`
	Templates []Template     `mapstructure:"templates,omitempty"`
	Contacts  []ContactPoint `mapstructure:"contacts,omitempty"`
}

func (c *Config) BuildServiceGraph() ([]*tree.Node[*Service], error) {
	return tree.NewForest(c.Services)
}
