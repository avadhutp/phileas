package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/avadhutp/phileas/command"
)

var (
	logger = log.WithFields(log.Fields{"package": "main"})
)

func main() {
	if err := command.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
