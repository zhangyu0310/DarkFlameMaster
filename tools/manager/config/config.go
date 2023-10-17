package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
)

// Config Configurations of server.
type Config struct {
	Server string `json:"server"`
}

var (
	globalConf atomic.Value
)

// ReadConfig reads the configuration from the given path and stores it in the given config.
func ReadConfig(path string) error {
	cfg := Config{}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Read config file failed, err:", err)
		return err
	}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Println("Unmarshal config file failed, err:", err)
		return err
	}
	if cfg.Server == "" {
		cfg.Server = "http://localhost:1219"
	}
	if !strings.HasPrefix(cfg.Server, "http://") {
		cfg.Server = "http://" + cfg.Server
	}
	StoreGlobalConfig(&cfg)
	return nil
}

// GetGlobalConfig returns the global configuration for this server.
// It should store configuration from command line and configuration file.
// Other parts of the system can read the global configuration use this function.
func GetGlobalConfig() *Config {
	return globalConf.Load().(*Config)
}

// StoreGlobalConfig stores a new config to the globalConf. It mostly uses in the test to avoid some data races.
func StoreGlobalConfig(config *Config) {
	globalConf.Store(config)
}
