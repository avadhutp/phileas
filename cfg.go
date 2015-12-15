package main

import (
	"fmt"

	"github.com/fatih/structs"
)

type common struct {
	Port int `ini:"port"`
}

// Cfg Config map for phileas.ini
type Cfg struct {
	Common common `ini:"common"`
}

func (cfg *Cfg) dump() {
	dumpSection(cfg.Common)
}

func dumpSection(sect interface{}) {
	m := structs.Map(sect)

	for name, val := range m {
		switch v := val.(type) {
		case int:
			logger.Info(fmt.Sprintf("%s = %d", name, v))
		default:
			logger.Info(fmt.Sprintf("%s = %s", name, v))
		}
	}
}
