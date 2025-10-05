package config

import (
	"github.com/fmotalleb/go-tools/tree"

	"github.com/go-playground/validator/v10"
)

type ServiceNode = *tree.Node[*Service]

type Config struct {
	Services  []*Service     `mapstructure:"services,omitempty" validate:"dive"`
	Templates []Template     `mapstructure:"templates,omitempty" validate:"dive"`
	Contacts  []ContactPoint `mapstructure:"contacts,omitempty" validate:"dive"`
}

func (c *Config) Validate() error {
	validate := validator.New(
		validator.WithRequiredStructEnabled(),
		validator.WithPrivateFieldValidation(),
	)
	_ = validate.RegisterValidation("workdir", workingDirValidator)
	return validate.Struct(c)
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
	if n.Data.Lazy && len(n.Children()) == 0 {
		return false
	}
	if !n.Data.Enabled {
		return false
	}
	return true
}

func weightServices(rootServices []ServiceNode) {
	for _, root := range rootServices {
		weightServices(root.Children())
		if root.Data.GetDependCount() == 0 {
			root.Traverse(IncreaseDependCount)
		}
	}
}
