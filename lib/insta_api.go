package lib

import (
	"time"

	"github.com/gedex/go-instagram/instagram"
	"github.com/jinzhu/gorm"
)

const (
	waitBetweenChecks = 10 * time.Hour
)

type InstaAPI struct {
	client *instagram.Client
	db     *gorm.DB
}

// NewInstaAPI Provider for InstaAPI
func NewInstaAPI(cfg *Cfg) *InstaAPI {
	i := new(InstaAPI)

	i.client = instagram.NewClient(nil)
	i.client.ClientID = cfg.Instagram.ClientID
	i.client.ClientSecret = cfg.Instagram.Secret
	i.client.AccessToken = cfg.Instagram.Token
	i.db = GetDB(cfg)

	return i
}

// SaveLikes Inserts instagram likes into the DB
func (i *InstaAPI) SaveLikes() {
	for {
		media, _, _ := i.client.Users.LikedMedia(nil)

		for _, m := range media {
			if i.isLocationOk(&m) {
				var e Entry
				i.db.FirstOrCreate(&e, Entry{
					Type:      "instagram",
					VendorID:  m.ID,
					Timestamp: m.CreatedTime,
				})
			}
		}

		time.Sleep(waitBetweenChecks)
	}
}

func (i *InstaAPI) isLocationOk(media *instagram.Media) bool {
	return media.Location != nil
}
