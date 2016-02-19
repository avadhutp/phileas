package command

import (
	"github.com/avadhutp/phileas/lib"

	"github.com/spf13/cobra"
)

var (
	esEnrichLocations = (*lib.EnrichmentService).EnrichLocation

	enrichLocationsCmd = &cobra.Command{
		Use:   "enrich-locations",
		Short: "Enrich locations",
		Long:  "Enrich locations with country and city information",
		Run:   enrichLocations,
	}
)

func enrichLocations(cmd *cobra.Command, args []string) {
	logger.Infof("Enriching Instagram likes for Phileas with city/country information; config's at %s", cfgPath)

	cfg := libNewCfg(cfgPath)
	db := libGetDB(cfg)

	enrichmentService := libNewEnrichmentService(cfg, db)
	esEnrichLocations(enrichmentService)

	logger.Info("Enrichment done!")
}
