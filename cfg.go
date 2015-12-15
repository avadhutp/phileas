package main

import (
	"fmt"

	"github.com/fatih/structs"
)

type common struct {
	Port string `ini:"port"`
}

type instagram struct {
	ClientID string `ini:"client_id"`
	Secret   string `ini:"secret"`
	Token    string `ini:"access_token"`
}

// Cfg Config map for phileas.ini
type Cfg struct {
	Common    common    `ini:"common"`
	Instagram instagram `ini:"instagram"`
}

func (cfg *Cfg) dump() {
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
