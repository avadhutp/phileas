package command

import (
	"fmt"

	"github.com/avadhutp/phileas/lib"
	"github.com/spf13/cobra"
)

var (
	backfillCmd = &cobra.Command{
		Use:   "backfill",
		Short: "Backfill the instagram likes",
		Long:  "Backfill the instagram likes; to be used when the service is first starting up",
		Run:   backfillPhileas,
	}
)

func backfillPhileas(cmd *cobra.Command, args []string) {
	logger.Info(fmt.Sprintf("Backfilling Instagram likes for Phileas; config's at %s", cfgPath))

	cfg := lib.NewCfg(cfgPath)
	db := lib.GetDB(cfg)

	instaAPI := lib.NewInstaAPI(cfg, db)
	instaAPI.Backfill("")

	logger.Info("Backfilling done!")
}
