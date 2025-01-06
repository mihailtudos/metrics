package server

import (
	"flag"
	"sync"
	"time"

	utils "github.com/mihailtudos/metrics/utils/configs"
)

type ServerConfig struct {
	Address       string
	StoreInterval time.Duration
	FileStorePath string
	Restore       bool
}

const (
	DefaultServerAddress = "localhost:8080"
	DefaultStoreInterval = 300
)

var (
	instance             *ServerConfig
	once                 sync.Once
	DefaultFileStorePath = "metrics-db.json"
)

func NewServerConfig() *ServerConfig {
	once.Do(func() {
		cfg := &struct {
			Address       string
			StoreInterval int
			FileStorePath string
			Restore       bool
		}{}

		flag.StringVar(&cfg.Address, "a", DefaultServerAddress, "server address")
		flag.IntVar(&cfg.StoreInterval, "i", DefaultStoreInterval, "metrics store interval")
		flag.StringVar(&cfg.FileStorePath, "f", DefaultFileStorePath, "file store path")
		flag.BoolVar(&cfg.Restore, "r", false, "restore metrics from file")

		flag.Parse()

		utils.OverrideStringEnvValueWithOsEnv(&cfg.Address, "ADDRESS")
		utils.OverrideStringEnvValueWithOsEnv(&cfg.StoreInterval, "STORE_INTERVAL")
		utils.OverrideStringEnvValueWithOsEnv(&cfg.FileStorePath, "FILE_STORAGE_PATH")
		utils.OverrideStringEnvValueWithOsEnv(&cfg.Restore, "RESTORE")

		serverConfig := &ServerConfig{
			Address:       cfg.Address,
			StoreInterval: time.Duration(cfg.StoreInterval) * time.Second,
			FileStorePath: cfg.FileStorePath,
			Restore:       cfg.Restore,
		}

		instance = serverConfig
	})

	return instance
}
