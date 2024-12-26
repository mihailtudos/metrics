package server

import (
	"flag"
	"sync"

	"github.com/mihailtudos/metrics/internal/infrastructure/config/utils"
)

type ServerConfig struct {
	Address string
}

const DefaultServerAddress = "localhost:8080"

var (
	instance *ServerConfig
	once     sync.Once
)

func NewServerConfig() *ServerConfig {
	once.Do(func() {
		cfg := &ServerConfig{}

		flag.StringVar(&cfg.Address, "a", DefaultServerAddress, "server address")

		flag.Parse()

		utils.OverrideStringEnvValueWithOsEnv(&cfg.Address, "ADDRESS")
		instance = cfg
	})

	return instance
}
