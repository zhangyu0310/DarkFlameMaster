package config

import (
	"encoding/json"
	"os"
	"strings"
	"sync/atomic"

	zlog "github.com/zhangyu0310/zlogger"
)

// Config Configurations of server.
type Config struct {
	// Database info
	DbType DbType `json:"DbType"`
	DbPath string `json:"DbPath"`
	// Seat info file
	SeatFileType SeatFileType `json:"SeatFileType"`
	SeatFile     string       `json:"SeatFile"`
	// Customer info file
	CustomerType       CustomerType       `json:"CustomerType"`
	CustomerFile       string             `json:"CustomerFile"`
	ChooseSeatStrategy ChooseSeatStrategy `json:"ChooseSeatStrategy"`
	LogPath            string             `json:"LogPath"`
	// Input label
	ProofName      string `json:"ProofName"`
	AdditionalName string `json:"AdditionalName"`
	// Admin
	RootUserName   string `json:"RootUserName"`
	ServicePort    uint   `json:"ServicePort"`
	AdminPort      uint   `json:"AdminPort"`
	AdminLocalOnly bool   `json:"AdminLocalOnly"`
}

var defaultConfig = Config{
	DbType:             LevelDB,
	DbPath:             "./run/db",
	SeatFileType:       JsonType,
	SeatFile:           "./data/奥斯卡长安国际影城-5号ALPD激光厅.json",
	CustomerType:       NoPay,
	CustomerFile:       "./data/customer.json",
	ChooseSeatStrategy: NoLimit,
	LogPath:            "./run/log",
	ProofName:          "",
	AdditionalName:     "",
	RootUserName:       "",
	ServicePort:        718,
	AdminPort:          1219,
	AdminLocalOnly:     true,
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
	QQNum  CustomerType = "qqnum"
	NoPay  CustomerType = "nopay"
	TestCF CustomerType = "test"
)

type ChooseSeatStrategy string

const (
	PayTimeOneByOne ChooseSeatStrategy = "paytimeonebyone"
	NoLimit         ChooseSeatStrategy = "nolimit"
)

var (
	globalConf atomic.Value
)

// InitializeConfig initialize the global config handler.
func InitializeConfig(enforceCmdArgs func(*Config)) {
	cfg := defaultConfig
	// Use command config cover config file.
	enforceCmdArgs(&cfg)
	StoreGlobalConfig(&cfg)
}

// ReadConfig reads the configuration from the given path and stores it in the given config.
func ReadConfig(path string, cfg *Config) {
	data, err := os.ReadFile(path)
	if err != nil {
		zlog.Fatal("Read config file failed, err:", err)
	}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		zlog.Fatal("Unmarshal config file failed, err:", err)
	}
	cfg.DbType = DbType(strings.ToLower(string(cfg.DbType)))
	cfg.SeatFileType = SeatFileType(strings.ToLower(string(cfg.SeatFileType)))
	cfg.CustomerType = CustomerType(strings.ToLower(string(cfg.CustomerType)))
	cfg.ChooseSeatStrategy = ChooseSeatStrategy(strings.ToLower(string(cfg.ChooseSeatStrategy)))
	if cfg.ServicePort == 0 {
		cfg.ServicePort = 718
	}
	if cfg.AdminPort == 0 {
		cfg.AdminPort = 1219
	}
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
