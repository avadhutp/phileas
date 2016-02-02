package command

import (
	"github.com/gin-gonic/gin"

	"github.com/avadhutp/phileas/lib"
	"github.com/spf13/cobra"
)

var (
	libNewService = lib.NewService
	serviceRun    = (*gin.Engine).Run
	startCmd      = &cobra.Command{
		Use:   "start",
		Short: "Start the web server",
		Long:  "Start the web server to provide Phileas as a service",
		Run:   startPhileas,
	}
)

func startPhileas(cmd *cobra.Command, args []string) {
	logger.Infof("Starting Phileas; config's at %s", cfgPath)

	cfg := libNewCfg(cfgPath)
	db := libGetDB(cfg)
	instaAPI := libNewInstaAPI(cfg, db)

	service := libNewService(cfg, db, instaAPI)
	serviceRun(service, ":"+cfg.Common.Port)
}
