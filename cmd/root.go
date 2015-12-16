package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cfgPath string
	logger  = log.WithFields(log.Fields{"package": "main"})
	// RootCmd This represents the binary itself
	RootCmd = &cobra.Command{
		Use:   "phileas",
		Short: "Phileas creates your bucket list",
		Long:  "Phileas creates your bucket list based on your instagram & reddit likes & upvotes, respectively",
		Run: func(cmd *cobra.Command, args []string) {
			// Empty root command
		},
	}
)

func init() {
	RootCmd.Flags().StringVarP(&cfgPath, "config", "c", "/etc/phileas.ini", "Absolute path to the config file. Refer to the online documentation at github.com/avadhutp/philease for more info")

	RootCmd.AddCommand(setupCmd)
	RootCmd.AddCommand(startCmd)
}
