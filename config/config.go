package config

type Config struct {
	Services []Service `mapstructure:"services,omitempty"`
}
