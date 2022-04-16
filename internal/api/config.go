package api

import "time"

type Config struct {
	HttpClientTimeout         time.Duration `mapstructure:"http-client-timeout"`
	HttpServerTimeout         time.Duration `mapstructure:"http-server-timeout"`
	HttpServerShutdownTimeout time.Duration `mapstructure:"http-server-shutdown-timeout"`
	BackendURL                []string      `mapstructure:"backend-url"`
	ConfigPath                string        `mapstructure:"config-path"`
	Host                      string        `mapstructure:"host"`
	Port                      string        `mapstructure:"port"`
	Hostname                  string        `mapstructure:"hostname"`
	H2C                       bool          `mapstructure:"h2c"`
	Unhealthy                 bool          `mapstructure:"unhealthy"`
	Unready                   bool          `mapstructure:"unready"`
	JWTSecret                 string        `mapstructure:"jwt-secret"`
}
