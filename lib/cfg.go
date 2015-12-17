package lib

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-ini/ini"
)

var (
	logger = log.WithFields(log.Fields{"file": "cfg.go"})
)

type common struct {
	Port        string `ini:"port"`
	MapquestKey string `ini:"mapquest_key"`
}

type instagramConfig struct {
	ClientID string `ini:"client_id"`
	Secret   string `ini:"secret"`
	Token    string `ini:"access_token"`
}

type mysqlConfig struct {
	Host     string `ini:"host"`
	Port     string `ini:"port"`
	Username string `ini:"username"`
	Password string `ini:"password"`
	Database string `ini:"schema"`
}

// Cfg Config map for phileas.ini
type Cfg struct {
	Common    common          `ini:"common"`
	Instagram instagramConfig `ini:"instagram"`
	Mysql     mysqlConfig     `ini:"mysql"`
}

// NewCfg Gets the config struct
func NewCfg(fileName string) *Cfg {
	var cfg Cfg
	err := ini.MapTo(&cfg, "phileas.ini")

	if err != nil {
		logger.Error("Cannot parse the config file: ", err)
	}

	return &cfg
}
