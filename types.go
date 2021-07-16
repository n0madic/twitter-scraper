package twitterscraper

import "time"

type (
	// Video type.
	Video struct {
		ID      string
		Preview string
		URL     string
	}

	// Tweet type.
	Tweet struct {
		Hashtags        []string
		HTML            string
		ID              string
		InReplyToStatus *Tweet
		IsQuoted        bool
		IsPin           bool
		IsReply         bool
		IsRetweet       bool
		Likes           int
		PermanentURL    string
		Photos          []string
		QuotedStatus    *Tweet
		Replies         int
		Retweets        int
		RetweetedStatus *Tweet
		Text            string
		TimeParsed      time.Time
		Timestamp       int64
		URLs            []string
		UserID          string
		Username        string
		Videos          []Video
	}

	// ProfileResult of scrapping.
	ProfileResult struct {
		Profile
		Error error
	}

	// TweetResult of scrapping.
	TweetResult struct {
		Tweet
		Error error
	}

	legacyUser struct {
		CreatedAt   string `json:"created_at"`
		Description string `json:"description"`
		Entities    struct {
			URL struct {
				Urls []struct {
					ExpandedURL string `json:"expanded_url"`
				} `json:"urls"`
			} `json:"url"`
		} `json:"entities"`
		FavouritesCount      int      `json:"favourites_count"`
		FollowersCount       int      `json:"followers_count"`
		FriendsCount         int      `json:"friends_count"`
		IDStr                string   `json:"id_str"`
		ListedCount          int      `json:"listed_count"`
		Name                 string   `json:"name"`
		Location             string   `json:"location"`
		PinnedTweetIdsStr    []string `json:"pinned_tweet_ids_str"`
		ProfileBannerURL     string   `json:"profile_banner_url"`
		ProfileImageURLHTTPS string   `json:"profile_image_url_https"`
		Protected            bool     `json:"protected"`
		ScreenName           string   `json:"screen_name"`
		StatusesCount        int      `json:"statuses_count"`
		Verified             bool     `json:"verified"`
	}

	fetchProfileFunc func(query string, maxProfilesNbr int, cursor string) ([]*Profile, string, error)
	fetchTweetFunc   func(query string, maxTweetsNbr int, cursor string) ([]*Tweet, string, error)
)
