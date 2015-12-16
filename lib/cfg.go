package lib

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/go-ini/ini"
)

var (
	logger = log.WithFields(log.Fields{"package": "lib"})
)

type common struct {
	Port string `ini:"port"`
}

type instagramConfig struct {
	ClientID string `ini:"client_id"`
	Secret   string `ini:"secret"`
	Token    string `ini:"access_token"`
}

// Cfg Config map for phileas.ini
type Cfg struct {
	Common    common          `ini:"common"`
	Instagram instagramConfig `ini:"instagram"`
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

// Dump Dumps the contents of the parsed config file, for debug purposes
func (cfg *Cfg) Dump() {
	dumpSection(cfg.Common, "Common")
	dumpSection(cfg.Instagram, "Instagram")
}

func dumpSection(sect interface{}, sectName string) {
	m := structs.Map(sect)

	for name, val := range m {
		switch v := val.(type) {
		case int:
			logger.Info(fmt.Sprintf("%s.%s = %d", sectName, name, v))
		default:
			logger.Info(fmt.Sprintf("%s.%s = %s", sectName, name, v))
		}
	}
}
