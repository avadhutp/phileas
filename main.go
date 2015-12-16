package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/avadhutp/phileas/lib"
)

var (
	logger = log.WithFields(log.Fields{"package": "main"})
)

func main() {
	cfg := lib.NewCfg("/etc/phileas.ini")
	cfg.Dump()

	service := lib.NewService(cfg)
	service.Run(":" + cfg.Common.Port)
}
