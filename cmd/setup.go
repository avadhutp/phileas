package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logger   = log.WithFields(log.Fields{"package": "main"})
	setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Setup Phileas, install DB, etc.",
		Long:  "Setup Phileas, install DB, make initial entries, etc. This command is idempotent.",
		Run:   setup,
	}
)

func init() {
	RootCmd.AddCommand(setupCmd)
}

func setup(cmd *cobra.Command, args []string) {
	logger.Info("Setting up Phileas...")
}
