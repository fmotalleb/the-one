package config

import "github.com/fmotalleb/go-tools/tree"

type Config struct {
	Services  []*Service     `mapstructure:"services,omitempty"`
	Templates []Template     `mapstructure:"templates,omitempty"`
	Contacts  []ContactPoint `mapstructure:"contacts,omitempty"`
}

func (c *Config) BuildServiceGraph() ([]*tree.Node[*Service], error) {
	// Config need to be sanitized at first
	// then create the tree
	// remove lazy nodes with no children
	// create abstract of runtime
	// abstract should be absolute and envfile should be loaded into env variables already
	// abstract will be passed to Boot service directly
	return tree.NewForest(c.Services)
}
