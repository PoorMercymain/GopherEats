// Package config populates server config data based on flags, environment variables or defaults
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
)

type serverConfig struct {
	HostGrpc    string `json:"grpc_address,omitempty"`
	DatabaseDSN string `env:"DATABASE_DSN" json:"database_dsn,omitempty"`
	ConfigFile  string `env:"CONFIG"`
}

// GetServerConfig returns server config params
func GetServerConfig() (config serverConfig) {

	var flagConfig serverConfig

	_ = env.Parse(&config)

	flagConfig = getServerFlags(config.ConfigFile)

	if flagConfig.HostGrpc != "" {
		config.HostGrpc = flagConfig.HostGrpc
	}

	if cfg.DatabaseDSN != "" {
		dbDSN = cfg.DatabaseDSN
	}

	return
}

func getServerFlags(configFile string) (config serverConfig) {

	if flag.Lookup("config") == nil && configFile == "" {
		fmt.Println("setting config file from flag")
		flag.StringVar(&configFile, "config", "", "config JSON file path")
	}

	if flag.Lookup("d") == nil {
		flag.StringVar(&dbDSN, "d", "", "database address")
	}
	flag.Parse()

	//sensible defaults to run in absence of flags and env vars
	configDefaults := &serverConfig{
		HostGrpc: "localhost:3200",
	}

	fmt.Println("config file ", configFile)

	if configFile != "" {
		fmt.Println("reading config file")
		dat, err := os.ReadFile(configFile)
		if err != nil {
			fmt.Println("failed to read config file %w", err)
		}
		fmt.Println(dat)
		err = json.Unmarshal(dat, configDefaults)
		if err != nil {
			fmt.Println("failed to unmarshal config file %w", err)
		}
		fmt.Println(configDefaults)

	}

	if gRPCHost == "" {
		gRPCHost = configDefaults.HostGrpc
	}

	if dbDSN == "" {
		dbDSN = configDefaults.DatabaseDSN
	}

	return
}
