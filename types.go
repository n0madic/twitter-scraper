package twitterscraper

import "time"

type (
	// Mention type.
	Mention struct {
		ID       string
		Username string
		Name     string
	}

	// Photo type.
	Photo struct {
		ID  string
		URL string
	}

	// Video type.
	Video struct {
		ID      string
		Preview string
		URL     string
	}

	// Tweet type.
	Tweet struct {
		Hashtags          []string
		HTML              string
		ID                string
		InReplyToStatus   *Tweet
		InReplyToStatusID string
		IsQuoted          bool
		IsPin             bool
		IsReply           bool
		IsRetweet         bool
		Likes             int
		Name              string
		Mentions          []Mention
		PermanentURL      string
		Photos            []Photo
		Place             *Place
		QuotedStatus      *Tweet
		QuotedStatusID    string
		Replies           int
		Retweets          int
		RetweetedStatus   *Tweet
		RetweetedStatusID string
		Text              string
		TimeParsed        time.Time
		Timestamp         int64
		URLs              []string
		UserID            string
		Username          string
		Videos            []Video
		Views             int
		SensitiveContent  bool
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

	legacyTweet struct {
		CreatedAt     string `json:"created_at"`
		FavoriteCount int    `json:"favorite_count"`
		FullText      string `json:"full_text"`
		Entities      struct {
			Hashtags []struct {
				Text string `json:"text"`
			} `json:"hashtags"`
			Media []struct {
				MediaURLHttps string `json:"media_url_https"`
				Type          string `json:"type"`
				URL           string `json:"url"`
			} `json:"media"`
			URLs []struct {
				ExpandedURL string `json:"expanded_url"`
				URL         string `json:"url"`
			} `json:"urls"`
			UserMentions []struct {
				IDStr      string `json:"id_str"`
				Name       string `json:"name"`
				ScreenName string `json:"screen_name"`
			} `json:"user_mentions"`
		} `json:"entities"`
		ExtendedEntities struct {
			Media []struct {
				IDStr                    string `json:"id_str"`
				MediaURLHttps            string `json:"media_url_https"`
				ExtSensitiveMediaWarning struct {
					AdultContent    bool `json:"adult_content"`
					GraphicViolence bool `json:"graphic_violence"`
					Other           bool `json:"other"`
				} `json:"ext_sensitive_media_warning"`
				Type      string `json:"type"`
				URL       string `json:"url"`
				VideoInfo struct {
					Variants []struct {
						Bitrate int    `json:"bitrate"`
						URL     string `json:"url"`
					} `json:"variants"`
				} `json:"video_info"`
			} `json:"media"`
		} `json:"extended_entities"`
		IDStr                 string `json:"id_str"`
		InReplyToStatusIDStr  string `json:"in_reply_to_status_id_str"`
		Place                 Place  `json:"place"`
		ReplyCount            int    `json:"reply_count"`
		RetweetCount          int    `json:"retweet_count"`
		RetweetedStatusIDStr  string `json:"retweeted_status_id_str"`
		RetweetedStatusResult struct {
			Result *result `json:"result"`
		} `json:"retweeted_status_result"`
		QuotedStatusIDStr string    `json:"quoted_status_id_str"`
		Time              time.Time `json:"time"`
		UserIDStr         string    `json:"user_id_str"`
		Views             struct {
			State string `json:"state"`
			Count string `json:"count"`
		} `json:"ext_views"`
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

	Place struct {
		ID          string `json:"id"`
		PlaceType   string `json:"place_type"`
		Name        string `json:"name"`
		FullName    string `json:"full_name"`
		CountryCode string `json:"country_code"`
		Country     string `json:"country"`
		BoundingBox struct {
			Type        string        `json:"type"`
			Coordinates [][][]float64 `json:"coordinates"`
		} `json:"bounding_box"`
	}

	fetchProfileFunc func(query string, maxProfilesNbr int, cursor string) ([]*Profile, string, error)
	fetchTweetFunc   func(query string, maxTweetsNbr int, cursor string) ([]*Tweet, string, error)
)
