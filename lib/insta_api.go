package lib

import (
	"github.com/gedex/go-instagram/instagram"
	"github.com/jinzhu/gorm"
	"net/url"
	"time"
)

const (
	waitBetweenChecks = 10 * time.Hour
	backfillWait      = 5 * time.Second
)

// DI
var (
	getLikedMedia = (*instagram.UsersService).LikedMedia
	timeSleep     = time.Sleep

	instaAPISaveLocation func(*InstaAPI, *instagram.Media) *Location
	instaAPISaveMedia    func(*InstaAPI, *instagram.Media)
	instaAPIBackfill     func(*InstaAPI, string)
	instaAPISaveLikes    func(*InstaAPI)
)

func init() {
	instaAPISaveLocation = (*InstaAPI).saveLocation
	instaAPISaveMedia = (*InstaAPI).saveMedia
	instaAPIBackfill = (*InstaAPI).Backfill
	instaAPISaveLikes = (*InstaAPI).SaveLikes
}

// InstaAPI encapsulate functionality for all instagram functionality
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
	media, _, _ := getLikedMedia(i.client.Users, nil)

	for _, m := range media {
		instaAPISaveMedia(i, &m)
	}

	timeSleep(waitBetweenChecks)
	instaAPISaveLikes(i)
}

// Backfill Puts in historical likes
func (i *InstaAPI) Backfill(maxLikeID string) {
	logger.Infof("Running backfill for %s", maxLikeID)

	media, after, _ := getLikedMedia(i.client.Users, &instagram.Parameters{MaxID: maxLikeID})
	afterURL, _ := url.Parse(after.NextURL)
	maxLikeID = afterURL.Query().Get("max_like_id")

	for _, m := range media {
		instaAPISaveMedia(i, &m)
	}

	if maxLikeID != "" {
		timeSleep(backfillWait)
		instaAPIBackfill(i, maxLikeID)
	}
}

func (i *InstaAPI) saveMedia(m *instagram.Media) {
	if !i.isLocationOk(m) {
		return
	}

	loc := instaAPISaveLocation(i, m)
	var e Entry
	i.db.FirstOrCreate(&e, Entry{
		Type:       "instagram",
		VendorID:   m.ID,
		Thumbnail:  m.Images.Thumbnail.URL,
		URL:        m.Link,
		Caption:    getCaption(m),
		Timestamp:  m.CreatedTime,
		LocationID: loc.ID,
	})
}

func (i *InstaAPI) saveLocation(m *instagram.Media) *Location {
	var l Location
	i.db.FirstOrCreate(&l, Location{
		Name: m.Location.Name,
		Lat:  m.Location.Latitude,
		Long: m.Location.Longitude,
	})

	return &l
}

func (i *InstaAPI) isLocationOk(media *instagram.Media) bool {
	return media.Location != nil
}

func getCaption(m *instagram.Media) string {
	if m.Caption != nil {
		return m.Caption.Text
	}

	return ""
}
