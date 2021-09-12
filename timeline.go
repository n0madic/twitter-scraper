package twitterscraper

import (
	"fmt"
	"strings"
	"time"
)

// timeline JSON object
type timeline struct {
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
								User struct {
									ID string `json:"id"`
								} `json:"user"`
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
			ReplaceEntry struct {
				Entry struct {
					Content struct {
						Operation struct {
							Cursor struct {
								Value      string `json:"value"`
								CursorType string `json:"cursorType"`
							} `json:"cursor"`
						} `json:"operation"`
					} `json:"content"`
				} `json:"entry"`
			} `json:"replaceEntry,omitempty"`
		} `json:"instructions"`
	} `json:"timeline"`
}

func (timeline *timeline) parseTweet(id string) *Tweet {
	if tweet, ok := timeline.GlobalObjects.Tweets[id]; ok {
		username := timeline.GlobalObjects.Users[tweet.UserIDStr].ScreenName
		tw := &Tweet{
			ID:           id,
			Likes:        tweet.FavoriteCount,
			PermanentURL: fmt.Sprintf("https://twitter.com/%s/status/%s", username, id),
			Replies:      tweet.ReplyCount,
			Retweets:     tweet.RetweetCount,
			Text:         tweet.FullText,
			UserID:       tweet.UserIDStr,
			Username:     username,
		}

		tm, err := time.Parse(time.RubyDate, tweet.CreatedAt)
		if err == nil {
			tw.TimeParsed = tm
			tw.Timestamp = tm.Unix()
		}

		if tweet.QuotedStatusIDStr != "" {
			tw.IsQuoted = true
			tw.QuotedStatus = timeline.parseTweet(tweet.QuotedStatusIDStr)
		}
		if tweet.InReplyToStatusIDStr != "" {
			tw.IsReply = true
			tw.InReplyToStatus = timeline.parseTweet(tweet.InReplyToStatusIDStr)
		}
		if tweet.RetweetedStatusIDStr != "" {
			tw.IsRetweet = true
			tw.RetweetedStatus = timeline.parseTweet(tweet.RetweetedStatusIDStr)
		}

		for _, pinned := range timeline.GlobalObjects.Users[tweet.UserIDStr].PinnedTweetIdsStr {
			if tweet.ConversationIDStr == pinned {
				tw.IsPin = true
				break
			}
		}

		for _, hash := range tweet.Entities.Hashtags {
			tw.Hashtags = append(tw.Hashtags, hash.Text)
		}
		for _, media := range tweet.Entities.Media {
			if media.Type == "photo" {
				tw.Photos = append(tw.Photos, media.MediaURLHttps)
			}
		}
		for _, media := range tweet.ExtendedEntities.Media {
			if media.Type == "video" {
				video := Video{
					ID:      media.IDStr,
					Preview: media.MediaURLHttps,
				}
				maxBitrate := 0
				for _, variant := range media.VideoInfo.Variants {
					if variant.Bitrate > maxBitrate {
						video.URL = strings.TrimSuffix(variant.URL, "?tag=10")
						maxBitrate = variant.Bitrate
					}
				}
				tw.Videos = append(tw.Videos, video)
			}
		}
		for _, url := range tweet.Entities.URLs {
			tw.URLs = append(tw.URLs, url.ExpandedURL)
		}

		tw.HTML = tweet.FullText
		tw.HTML = reHashtag.ReplaceAllStringFunc(tw.HTML, func(hashtag string) string {
			return fmt.Sprintf(`<a href="https://twitter.com/hashtag/%s">%s</a>`,
				strings.TrimPrefix(hashtag, "#"),
				hashtag,
			)
		})
		tw.HTML = reUsername.ReplaceAllStringFunc(tw.HTML, func(username string) string {
			return fmt.Sprintf(`<a href="https://twitter.com/%s">%s</a>`,
				strings.TrimPrefix(username, "@"),
				username,
			)
		})
		tw.HTML = reTwitterURL.ReplaceAllStringFunc(tw.HTML, func(tco string) string {
			for _, entity := range tweet.Entities.URLs {
				if tco == entity.URL {
					return fmt.Sprintf(`<a href="%s">%s</a>`, entity.ExpandedURL, tco)
				}
			}
			for _, entity := range tweet.Entities.Media {
				if tco == entity.URL {
					return fmt.Sprintf(`<br><a href="%s"><img src="%s"/></a>`, tco, entity.MediaURLHttps)
				}
			}
			return tco
		})
		tw.HTML = strings.Replace(tw.HTML, "\n", "<br>", -1)
		return tw
	}
	return nil
}

func (timeline *timeline) parseTweets() ([]*Tweet, string) {
	var cursor string
	var pinnedTweet *Tweet
	var orderedTweets []*Tweet
	for _, instruction := range timeline.Timeline.Instructions {
		if instruction.PinEntry.Entry.Content.Item.Content.Tweet.ID != "" {
			if tweet := timeline.parseTweet(instruction.PinEntry.Entry.Content.Item.Content.Tweet.ID); tweet != nil {
				pinnedTweet = tweet
			}
		}
		for _, entry := range instruction.AddEntries.Entries {
			if tweet := timeline.parseTweet(entry.Content.Item.Content.Tweet.ID); tweet != nil {
				orderedTweets = append(orderedTweets, tweet)
			}
			if entry.Content.Operation.Cursor.CursorType == "Bottom" {
				cursor = entry.Content.Operation.Cursor.Value
			}
		}
		if instruction.ReplaceEntry.Entry.Content.Operation.Cursor.CursorType == "Bottom" {
			cursor = instruction.ReplaceEntry.Entry.Content.Operation.Cursor.Value
		}
	}
	if pinnedTweet != nil && len(orderedTweets) > 0 {
		orderedTweets = append([]*Tweet{pinnedTweet}, orderedTweets...)
	}
	return orderedTweets, cursor
}

func (timeline *timeline) parseUsers() ([]*Profile, string) {
	users := make(map[string]Profile)

	for id, user := range timeline.GlobalObjects.Users {
		users[id] = parseProfile(user)
	}

	var cursor string
	var orderedProfiles []*Profile
	for _, instruction := range timeline.Timeline.Instructions {
		for _, entry := range instruction.AddEntries.Entries {
			if profile, ok := users[entry.Content.Item.Content.User.ID]; ok {
				orderedProfiles = append(orderedProfiles, &profile)
			}
			if entry.Content.Operation.Cursor.CursorType == "Bottom" {
				cursor = entry.Content.Operation.Cursor.Value
			}
		}
		if instruction.ReplaceEntry.Entry.Content.Operation.Cursor.CursorType == "Bottom" {
			cursor = instruction.ReplaceEntry.Entry.Content.Operation.Cursor.Value
		}
	}
	return orderedProfiles, cursor
}
