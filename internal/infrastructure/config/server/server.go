package server

import (
	"flag"
	"sync"
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

		instance = cfg
	})

	return instance
}
