package command

import (
	"github.com/avadhutp/phileas/lib"

	"github.com/spf13/cobra"
)

var (
	enrichCmd = &cobra.Command{
		Use:   "enrich",
		Short: "Enrich locations",
		Long:  "Enrich locations with country and city information",
		Run:   enrichPhileas,
	}
)

func enrichPhileas(cmd *cobra.Command, args []string) {
	logger.Infof("Enriching Instagram likes for Phileas with city/country information; config's at %s", cfgPath)

	cfg := libNewCfg(cfgPath)
	db := libGetDB(cfg)

	enrichmentService := lib.NewEnrichmentService(cfg, db)
	enrichmentService.EnrichLocation()

	logger.Info("Enrichment done!")
}
