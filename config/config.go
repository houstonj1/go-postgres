package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config shape
type Config struct {
	DBUsername string `envconfig:"DB_USERNAME" default:"postgres"`
	DBPassword string `envconfig:"DB_PASSWORD" default:""`
}

// NewConfig return new Config
func NewConfig() *Config {
	config := &Config{}
	envconfig.MustProcess("", config)

	return config
}

// Print return Config as string
func (c *Config) Print() string {
	return fmt.Sprintf("%#v", c)
}
