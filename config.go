package swim

import (
	"time"

	"github.com/ganmacs/swim/logger"
)

type Config struct {
	BindAddr string
	BindPort int

	ProbeInterval time.Duration
	ProbeTimeout  time.Duration

	SwimIterval   time.Duration
	SwimNodeCount int

	transport *Transport

	Logger *logger.Logger
}

func DefaultConfig() *Config {
	return &Config{
		BindAddr:      "0.0.0.0",
		BindPort:      8000,
		ProbeInterval: 2,
		ProbeTimeout:  time.Second * 3,
		SwimIterval:   time.Second * 2,
		SwimNodeCount: 3,
	}
}
