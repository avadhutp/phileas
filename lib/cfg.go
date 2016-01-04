package lib

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-ini/ini"
)

var (
	logger   = log.WithFields(log.Fields{"file": "cfg.go"})
	logErr   = logger.Error
	iniMapTo = ini.MapTo
)

type common struct {
	Port          string `ini:"port"`
	MapquestKey   string `ini:"mapquest_key"`
	GoogleMapsKey string `ini:"google_maps_key"`
}

type instagramConfig struct {
	ClientID string `ini:"client_id"`
	Secret   string `ini:"secret"`
	Token    string `ini:"access_token"`
}

type mysqlConfig struct {
	Host     string `ini:"host"`
	Port     string `ini:"port"`
	Username string `ini:"username"`
	Password string `ini:"password"`
	Database string `ini:"schema"`
}

type yelpConfig struct {
	ConsumerKey       string `ini:"consumer_key"`
	ConsumerSecret    string `ini:"consumer_secret"`
	AccessToken       string `ini:"access_token"`
	AccessTokenSecret string `ini:"access_token_secret"`
}

// Cfg Config map for phileas.ini
type Cfg struct {
	Common    common          `ini:"common"`
	Instagram instagramConfig `ini:"instagram"`
	Mysql     mysqlConfig     `ini:"mysql"`
	Yelp      yelpConfig      `ini:"yelp"`
}

// NewCfg Gets the config struct
func NewCfg(fileName string) *Cfg {
	var cfg Cfg
	err := iniMapTo(&cfg, "phileas.ini")

	if err != nil {
		logErr("Cannot parse the config file: ", err)
		return nil
	}

	return &cfg
}
