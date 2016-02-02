package command

import (
	"github.com/avadhutp/phileas/lib"

	"github.com/spf13/cobra"
)

var (
	instaAPIBackfill = (*lib.InstaAPI).Backfill
	backfillCmd      = &cobra.Command{
		Use:   "backfill",
		Short: "Backfill the instagram likes",
		Long:  "Backfill the instagram likes; to be used when the service is first starting up",
		Run:   backfillPhileas,
	}
)

func backfillPhileas(cmd *cobra.Command, args []string) {
	logger.Infof("Backfilling Instagram likes for Phileas; config's at %s", cfgPath)

	cfg := libNewCfg(cfgPath)
	db := libGetDB(cfg)

	instaAPI := libNewInstaAPI(cfg, db)
	instaAPIBackfill(instaAPI, "")

	logger.Info("Backfilling done!")
}
