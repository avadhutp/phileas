package lib

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gedex/go-instagram/instagram"
	"github.com/jinzhu/gorm"
)

const (
	waitBetweenChecks = 10 * time.Hour
	backfillWait      = 5 * time.Second
)

type InstaAPI struct {
	client *instagram.Client
	db     *gorm.DB
}

// NewInstaAPI Provider for InstaAPI
func NewInstaAPI(cfg *Cfg, db *gorm.DB) *InstaAPI {
	i := new(InstaAPI)

	i.client = instagram.NewClient(nil)
	i.client.ClientID = cfg.Instagram.ClientID
	i.client.ClientSecret = cfg.Instagram.Secret
	i.client.AccessToken = cfg.Instagram.Token
	i.db = db

	return i
}

// SaveLikes Inserts instagram likes into the DB
func (i *InstaAPI) SaveLikes() {
	for {
		media, _, _ := i.client.Users.LikedMedia(nil)

		for _, m := range media {
			i.saveMedia(&m)
		}

		time.Sleep(waitBetweenChecks)
	}
}

// Backfill Puts in historical likes
func (i *InstaAPI) Backfill(maxLikeID string) {
	media, after, _ := i.client.Users.LikedMedia(&instagram.Parameters{MaxID: maxLikeID})
	afterURL, _ := url.Parse(after.NextURL)
	maxLikeID = afterURL.Query().Get("max_like_id")

	for _, m := range media {
		i.saveMedia(&m)
	}

	if maxLikeID != "" {
		time.Sleep(backfillWait)
		i.Backfill(maxLikeID)
	}
}

// MediaInfo Retrieves info about one one instagram post
func (i *InstaAPI) MediaInfo(mediaId string) *instagram.Media {
	if info, err := i.client.Media.Get(mediaId); err == nil {
		return info
	} else {
		logger.Error(fmt.Sprintf("Cannot fetch media info: %s", err.Error()))
		return nil
	}
}

func (i *InstaAPI) saveMedia(m *instagram.Media) {
	if i.isLocationOk(m) {
		loc := i.saveLocation(m)
		var e Entry
		i.db.FirstOrCreate(&e, Entry{
			Type:       "instagram",
			VendorID:   m.ID,
			Timestamp:  m.CreatedTime,
			LocationID: loc.ID,
		})
	}
}

func (i *InstaAPI) saveLocation(m *instagram.Media) *Location {
	var l Location
	i.db.FirstOrCreate(&l, Location{
		Name:      m.Location.Name,
		Lat:       m.Location.Latitude,
		Long:      m.Location.Longitude,
		Thumbnail: m.Images.Thumbnail.URL,
		URL:       m.Link,
		Caption:   "test caption",
	})

	return &l
}

func (i *InstaAPI) isLocationOk(media *instagram.Media) bool {
	return media.Location != nil
}
