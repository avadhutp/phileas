package lib

import (
	"fmt"

	"github.com/gedex/go-instagram/instagram"
)

type InstaAPI struct {
	client *instagram.Client
}

// NewInstaAPI Provider for InstaAPI
func NewInstaAPI(cfg *Cfg) *InstaAPI {
	i := new(InstaAPI)

	i.client = instagram.NewClient(nil)
	i.client.ClientID = cfg.Instagram.ClientID
	i.client.ClientSecret = cfg.Instagram.Secret
	i.client.AccessToken = cfg.Instagram.Token

	return i
}

func (i *InstaAPI) SaveLikes() {
	media, _, _ := i.client.Users.LikedMedia(nil)

	for _, m := range media {
		if i.isLocationOk(&m) {
			logger.Info(fmt.Sprintf("ID: %s; Location: %s", m.ID, m.Location.Name))
		}
	}
}

func (i *InstaAPI) isLocationOk(media *instagram.Media) bool {
	return media.Location != nil
}
