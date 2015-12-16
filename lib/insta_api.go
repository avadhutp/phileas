package lib

import (
	"github.com/gedex/go-instagram/instagram"
)

var (
	instaAPI *instagram.Client
)

func initInstaAPI(cfg *Cfg) {
	instaAPI = instagram.NewClient(nil)

	instaAPI.ClientID = cfg.Instagram.ClientID
	instaAPI.ClientSecret = cfg.Instagram.Secret
	instaAPI.AccessToken = cfg.Instagram.Token
}
