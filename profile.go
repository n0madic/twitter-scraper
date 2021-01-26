package twitterscraper

import (
	"fmt"
	"net/http"
	"time"
)

// Profile of twitter user.
type Profile struct {
	Avatar         string
	Banner         string
	Biography      string
	Birthday       string
	FollowersCount int
	FollowingCount int
	FriendsCount   int
	IsPrivate      bool
	IsVerified     bool
	Joined         *time.Time
	LikesCount     int
	ListedCount    int
	Location       string
	Name           string
	PinnedTweetIDs []string
	TweetsCount    int
	URL            string
	UserID         string
	Username       string
	Website        string
}

// GetProfile return parsed user profile.
func (s *Scraper) GetProfile(username string) (Profile, error) {
	var jsn user
	req, err := http.NewRequest("GET", "https://api.twitter.com/graphql/4S2ihIKfF3xhp-ENxvUAfQ/UserByScreenName?variables=%7B%22screen_name%22%3A%22"+username+"%22%2C%22withHighlightedLabel%22%3Atrue%7D", nil)
	if err != nil {
		return Profile{}, err
	}

	err = s.RequestAPI(req, &jsn)
	if err != nil {
		return Profile{}, err
	}

	if len(jsn.Errors) > 0 {
		return Profile{}, fmt.Errorf("%s", jsn.Errors[0].Message)
	}

	if jsn.Data.User.RestID == "" {
		return Profile{}, fmt.Errorf("rest_id not found")
	}

	if jsn.Data.User.Legacy.ScreenName == "" {
		return Profile{}, fmt.Errorf("either @%s does not exist or is private", username)
	}

	user := jsn.Data.User.Legacy

	profile := Profile{
		Avatar:         user.ProfileImageURLHTTPS,
		Banner:         user.ProfileBannerURL,
		Biography:      user.Description,
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FavouritesCount,
		FriendsCount:   user.FriendsCount,
		IsPrivate:      user.Protected,
		IsVerified:     user.Verified,
		LikesCount:     user.FavouritesCount,
		ListedCount:    user.ListedCount,
		Location:       user.Location,
		Name:           user.Name,
		PinnedTweetIDs: user.PinnedTweetIdsStr,
		TweetsCount:    user.StatusesCount,
		URL:            "https://twitter.com/" + user.ScreenName,
		UserID:         jsn.Data.User.RestID,
		Username:       user.ScreenName,
	}

	tm, err := time.Parse(time.RubyDate, user.CreatedAt)
	if err == nil {
		tm = tm.UTC()
		profile.Joined = &tm
	}

	if len(user.Entities.URL.Urls) > 0 {
		profile.Website = user.Entities.URL.Urls[0].ExpandedURL
	}

	return profile, nil
}

// GetProfile wrapper for default scraper
func GetProfile(username string) (Profile, error) {
	return defaultScraper.GetProfile(username)
}
