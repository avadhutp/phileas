package command

import (
	log "github.com/Sirupsen/logrus"
	"github.com/avadhutp/phileas/lib"
	"github.com/spf13/cobra"
)

var (
	cfgPath                 string
	logger                  = log.WithFields(log.Fields{"package": "command"})
	libNewCfg               = lib.NewCfg
	libGetDB                = lib.GetDB
	libNewInstaAPI          = lib.NewInstaAPI
	libNewEnrichmentService = lib.NewEnrichmentService

	// RootCmd This represents the binary itself
	RootCmd = &cobra.Command{
		Use:   "phileas",
		Short: "Phileas creates your bucket list",
		Long:  "Phileas creates your bucket list based on your instagram & reddit likes & upvotes, respectively",
		Run:   nil,
	}
)

func init() {
	RootCmd.Flags().StringVarP(&cfgPath, "config", "c", "/etc/phileas.ini", "Absolute path to the config file. Refer to the online documentation at github.com/avadhutp/philease for more info")

	RootCmd.AddCommand(setupCmd)
	RootCmd.AddCommand(backfillCmd)
	RootCmd.AddCommand(startCmd)
	RootCmd.AddCommand(enrichCmd)
}
