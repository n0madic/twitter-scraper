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
		Hashtags     []string
		HTML         string
		ID           string
		IsQuoted     bool
		IsPin        bool
		IsReply      bool
		IsRetweet    bool
		Likes        int
		PermanentURL string
		Photos       []string
		Replies      int
		Retweets     int
		Text         string
		TimeParsed   time.Time
		Timestamp    int64
		URLs         []string
		UserID       string
		Username     string
		Videos       []Video
	}

	// Result of scrapping.
	Result struct {
		Tweet
		Error error
	}

	// timeline JSON
	timeline struct {
		GlobalObjects struct {
			Tweets map[string]struct {
				ConversationIDStr string `json:"conversation_id_str"`
				CreatedAt         string `json:"created_at"`
				FavoriteCount     int    `json:"favorite_count"`
				FullText          string `json:"full_text"`
				Entities          struct {
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
				} `json:"entities"`
				ExtendedEntities struct {
					Media []struct {
						IDStr         string `json:"id_str"`
						MediaURLHttps string `json:"media_url_https"`
						Type          string `json:"type"`
						VideoInfo     struct {
							Variants []struct {
								Bitrate int    `json:"bitrate,omitempty"`
								URL     string `json:"url"`
							} `json:"variants"`
						} `json:"video_info"`
					} `json:"media"`
				} `json:"extended_entities"`
				InReplyToStatusIDStr string    `json:"in_reply_to_status_id_str"`
				ReplyCount           int       `json:"reply_count"`
				RetweetCount         int       `json:"retweet_count"`
				RetweetedStatusIDStr string    `json:"retweeted_status_id_str"`
				QuotedStatusIDStr    string    `json:"quoted_status_id_str"`
				Time                 time.Time `json:"time"`
				UserIDStr            string    `json:"user_id_str"`
			} `json:"tweets"`
			Users map[string]struct {
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
			} `json:"users"`
		} `json:"globalObjects"`
		Timeline struct {
			Instructions []struct {
				AddEntries struct {
					Entries []struct {
						Content struct {
							Item struct {
								Content struct {
									Tweet struct {
										ID string `json:"id"`
									} `json:"tweet"`
								} `json:"content"`
							} `json:"item"`
							Operation struct {
								Cursor struct {
									Value      string `json:"value"`
									CursorType string `json:"cursorType"`
								} `json:"cursor"`
							} `json:"operation"`
							TimelineModule struct {
								Items []struct {
									Item struct {
										ClientEventInfo struct {
											Details struct {
												GuideDetails struct {
													TransparentGuideDetails struct {
														TrendMetadata struct {
															TrendName string `json:"trendName"`
														} `json:"trendMetadata"`
													} `json:"transparentGuideDetails"`
												} `json:"guideDetails"`
											} `json:"details"`
										} `json:"clientEventInfo"`
									} `json:"item"`
								} `json:"items"`
							} `json:"timelineModule"`
						} `json:"content,omitempty"`
					} `json:"entries"`
				} `json:"addEntries"`
				PinEntry struct {
					Entry struct {
						Content struct {
							Item struct {
								Content struct {
									Tweet struct {
										ID string `json:"id"`
									} `json:"tweet"`
								} `json:"content"`
							} `json:"item"`
						} `json:"content"`
					} `json:"entry"`
				} `json:"pinEntry,omitempty"`
			} `json:"instructions"`
		} `json:"timeline"`
	}

	fetchFunc func(user string, maxTweetsNbr int, cursor string) ([]*Tweet, string, error)
)
