package cmd

import (
	"fmt"

	"github.com/avadhutp/phileas/lib"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
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
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", cfg.Mysql.Username, cfg.Mysql.Password, cfg.Mysql.Host, cfg.Mysql.Port, cfg.Mysql.Database)

	if db, err := gorm.Open("mysql", connStr); err != nil {
		panic(err)
	} else {
		db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&lib.Entry{})
	}

	logger.Info("Phileas is ready!")
}
