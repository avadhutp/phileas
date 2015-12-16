package cmd

import (
	"fmt"

	"github.com/avadhutp/phileas/lib"

	"github.com/spf13/cobra"
)

var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the web server",
		Long:  "Start the web server to provide Phileas as a service",
		Run:   startPhileas,
	}
)

func startPhileas(cmd *cobra.Command, args []string) {
	logger.Info(fmt.Sprintf("Starting Phileas; config at %s", cfgPath))

	cfg := lib.NewCfg(cfgPath)
	service := lib.NewService(cfg)

	service.Run(":" + cfg.Common.Port)
}
