package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Setup Phileas, install DB, etc.",
		Long:  "Setup Phileas, install DB, make initial entries, etc. This command is idempotent.",
		Run:   setup,
	}
)

func setup(cmd *cobra.Command, args []string) {
	logger.Info(fmt.Sprintf("Setting up Phileas; config at %s", cfgPath))
}
