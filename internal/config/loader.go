package config

import (
	"github.com/kelseyhightower/envconfig"
)

func Load() (*Config, error) {
	settings := &Config{}

	err := envconfig.Process("", settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}
