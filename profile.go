package twitterscraper

import (
	"fmt"
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
func GetProfile(username string) (Profile, error) {
	userID, err := GetUserIDByScreenName(username)
	if err != nil {
		return Profile{}, err
	}

	req, err := newRequest("GET", "https://twitter.com/i/api/2/timeline/profile/"+userID+".json")
	if err != nil {
		return Profile{}, err
	}

	q := req.URL.Query()
	q.Add("count", "20")
	q.Add("userId", userID)
	req.URL.RawQuery = q.Encode()

	var timeline timeline
	err = requestAPI(req, &timeline)
	if err != nil {
		return Profile{}, err
	}

	user, found := timeline.GlobalObjects.Users[userID]
	if !found {
		return Profile{}, fmt.Errorf("either @%s does not exist or is private", username)
	}

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
		UserID:         user.IDStr,
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
