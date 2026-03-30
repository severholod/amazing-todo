package core_http_server

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	Address         string        `envconfig:"ADDRESS" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" required:"true"`
}

func NewConfig() (Config, error) {
	var conf Config
	if err := envconfig.Process("HTTP", &conf); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}
	return conf, nil
}

func NewConfigMust() Config {
	conf, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get HTTP server config: %w", err)
		panic(err)
	}
	return conf
}
