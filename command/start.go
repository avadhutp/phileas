package command

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
	db := lib.GetDB(cfg)

	instaAPI := lib.NewInstaAPI(cfg, db)
	// go instaAPI.Backfill("")

	// enrichment := lib.NewEnrichmentService(cfg, db)
	// go enrichment.EnrichLocation()
	// go enrichment.EnrichYelp()

	service := lib.NewService(cfg, db, instaAPI)
	service.Run(":" + cfg.Common.Port)
}
