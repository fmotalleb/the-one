package config

import "github.com/fmotalleb/the-one/logging"

var log = logging.LazyLogger("core.config")

type Config struct {
	Services []Service      `mapstructure:"services,omitempty"`
	Contacts []ContactPoint `mapstructure:"contacts,omitempty"`
}
