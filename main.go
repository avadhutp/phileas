package main

import (
	"github.com/avadhutp/phileas/cmd"

	log "github.com/Sirupsen/logrus"
)

var (
	logger = log.WithFields(log.Fields{"package": "main"})
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
