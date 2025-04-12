package config

type Config struct {
	Includes []string  `mapstructure:"include,omitempty"`
	Services []Service `mapstructure:"services,omitempty"`
}
