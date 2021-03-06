package command

import (
	"github.com/avadhutp/phileas/lib"
	// mysql import is unnamed for use with gorm
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

var (
	dbSet         = (*gorm.DB).Set
	dbAutoMigrate = (*gorm.DB).AutoMigrate

	setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Setup Phileas, install DB, etc.",
		Long:  "Setup Phileas, install DB, make initial entries, etc. This command is idempotent.",
		Run:   setup,
	}
)

func setup(cmd *cobra.Command, args []string) {
	logger.Infof("Setting up Phileas; config at %s", cfgPath)

	cfg := libNewCfg(cfgPath)

	db := libGetDB(cfg)
	dbSet(db, "gorm:table_options", "ENGINE=InnoDB")
	dbAutoMigrate(db, &lib.Entry{}, &lib.Location{})

	logger.Info("Phileas is ready!")
}
