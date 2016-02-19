package command

import (
	"github.com/avadhutp/phileas/lib"

	"github.com/spf13/cobra"
)

const (
	cmdDesc = `
	Enrich the locations table with additional information
	If you specify "enrich locations", Phileas will fetch country and city information for all the locations
	If you specify "enrich google-places", Phileas will fetch the place IDs for all locations; this will enable you to bookmark them to your Google account
	`
)

var (
	esEnrichLocations      = (*lib.EnrichmentService).EnrichLocation
	esEnrichGooglePlaceIDs = (*lib.EnrichmentService).EnrichGooglePlacesIDs

	enrichCmd = &cobra.Command{
		Use:   "enrich [locations|google-places]",
		Short: "Enrich locations or google places IDs",
		Long:  cmdDesc,
		Run:   enrich,
	}
)

func enrich(cmd *cobra.Command, args []string) {
	enrichmentType := args[0]

	logger.Infof("Enriching Instagram likes for %s; config's at %s", enrichmentType, cfgPath)

	cfg := libNewCfg(cfgPath)
	db := libGetDB(cfg)

	enrichmentService := libNewEnrichmentService(cfg, db)

	if enrichmentType == "locations" {
		esEnrichLocations(enrichmentService)
	} else if enrichmentType == "google-places" {
		esEnrichGooglePlaceIDs(enrichmentService)
	} else {
		logger.Errorf("Unknown enrichment type supplied as an argument: %s", enrichmentType)
	}

	logger.Info("Enrichment done!")
}
