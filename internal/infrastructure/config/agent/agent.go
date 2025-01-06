package agent

import (
	"flag"
	"fmt"
	"sync"
	"time"

	utils "github.com/mihailtudos/metrics/utils/configs"
)

type AgentConfig struct {
	ServerAddress  string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

const (
	DefaultServerAddress  = "localhost:8080"
	DefaultReportInterval = 10
	DefaultPollInterval   = 2
)

var (
	instance *AgentConfig
	once     sync.Once
)

func NewAgentConfig() *AgentConfig {
	once.Do(func() {
		type cfgStr struct {
			ServerAddress  string
			ReportInterval int
			PollInterval   int
		}

		flagsStr := cfgStr{}

		flag.StringVar(&flagsStr.ServerAddress, "a", DefaultServerAddress, "server address")
		flag.IntVar(&flagsStr.ReportInterval, "r", DefaultReportInterval, "report interval in seconds")
		flag.IntVar(&flagsStr.PollInterval, "p", DefaultPollInterval, "poll interval in seconds")

		flag.Parse()

		utils.OverrideStringEnvValueWithOsEnv(&flagsStr.ServerAddress, "ADDRESS")

		if err := utils.OverrideIntEnvValueWithOsEnv(&flagsStr.ReportInterval, "REPORT_INTERVAL"); err != nil {
			panic(err)
		}

		if err := utils.OverrideIntEnvValueWithOsEnv(&flagsStr.PollInterval, "POLL_INTERVAL"); err != nil {
			panic(err)
		}

		instance = &AgentConfig{
			ServerAddress:  fmt.Sprintf("http://%s", flagsStr.ServerAddress),
			ReportInterval: time.Duration(flagsStr.ReportInterval) * time.Second,
			PollInterval:   time.Duration(flagsStr.PollInterval) * time.Second,
		}
	})

	return instance
}
