package config

import "sync/atomic"

// Config Configurations of server.
type Config struct {
	// Database info
	DbType DbType
	DbPath string
	// Seat info file
	SeatFileType SeatFileType
	SeatFile     string
	// Customer info file
	CustomerType       CustomerType
	CustomerFile       string
	ChooseSeatStrategy ChooseSeatStrategy
}

type SeatFileType string

const (
	JsonType SeatFileType = "json"
	TestType SeatFileType = "test"
)

type DbType string

const (
	LevelDB DbType = "leveldb"
	MySQL   DbType = "mysql"
	ZDB     DbType = "zdb"
)

type CustomerType string

const (
	AliPay CustomerType = "alipay"
	WeChat CustomerType = "wechat"
	NoPay  CustomerType = "nopay"
	TestCF CustomerType = "test"
)

type ChooseSeatStrategy string

const (
	PayTimeOneByOne ChooseSeatStrategy = "ptonebyone"
	NoLimit         ChooseSeatStrategy = "nolimit"
)

var (
	globalConf atomic.Value
)

// InitializeConfig initialize the global config handler.
func InitializeConfig(enforceCmdArgs func(*Config)) {
	cfg := Config{}
	// Use command config cover config file.
	enforceCmdArgs(&cfg)
	StoreGlobalConfig(&cfg)
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
