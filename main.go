package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/avadhutp/phileas/cmd"
)

var (
	logger = log.WithFields(log.Fields{"package": "main"})
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
