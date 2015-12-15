package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-ini/ini"
)

var (
	logger = log.WithFields(log.Fields{"package": "main"})
)

func main() {
	cfg := getConfig()
	cfg.dump()

	service := getService()
	service.Run(":" + cfg.Common.Port)
}

func getConfig() *Cfg {
	var cfg Cfg
	err := ini.MapTo(&cfg, "phileas.ini")

	if err != nil {
		logger.Error("Cannot parse the config file: ", err)
	}

	return &cfg
}
