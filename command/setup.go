package command

import (
	"fmt"

	"github.com/avadhutp/phileas/lib"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

var (
	setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Setup Phileas, install DB, etc.",
		Long:  "Setup Phileas, install DB, make initial entries, etc. This command is idempotent.",
		Run:   setup,
	}
)

func setup(cmd *cobra.Command, args []string) {
	logger.Info(fmt.Sprintf("Setting up Phileas; config at %s", cfgPath))

	cfg := lib.NewCfg(cfgPath)

	db := lib.GetDB(cfg)
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&lib.Location{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&lib.Entry{})

	logger.Info("Phileas is ready!")
}
