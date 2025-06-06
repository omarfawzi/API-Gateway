package config

type Config struct {
	Cluster     string `envconfig:"APP_CLUSTER" default:"dev"`
	Environment string `envconfig:"APP_ENVIRONMENT" default:"local"`
	Port        int    `envconfig:"APP_PORT" default:"8080"`
	Sentry      SentryConfig
	Logger      LoggerConfig
}

type LoggerConfig struct {
	Level  string `envconfig:"LOGGER_LEVEL" default:"INFO"`
	Enable bool   `envconfig:"LOGGER_ENABLE" default:"false"`
}

type SentryConfig struct {
	Enable bool   `envconfig:"SENTRY_ENABLE" default:"false"`
	Dsn    string `envconfig:"SENTRY_DSN"`
	Debug  bool   `envconfig:"SENTRY_DEBUG" default:"false"`
}
