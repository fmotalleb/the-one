package config

import (
	"github.com/fmotalleb/go-tools/tree"
)

type ServiceNode = *tree.Node[*Service]

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
	rootServices, err := tree.NewForest(c.Services)
	if err != nil {
		return nil, err
	}
	rootServices = tree.ShakeForest(rootServices, treeFilter)
	weightServices(rootServices)
	return rootServices, err
}

func treeFilter(n ServiceNode) bool {
	// Remove lazy nodes with no child
	if n.Data.Lazy.UnwrapOr(false) && len(n.Children()) == 0 {
		return false
	}
	return false
}

func weightServices(rootServices []ServiceNode) {
	for _, root := range rootServices {
		weightServices(root.Children())
		if root.Data.GetDependCount() == 0 {
			root.Traverse(IncreaseDependCount)
		}
	}
}
