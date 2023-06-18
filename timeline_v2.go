package twitterscraper

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type result struct {
	Typename string `json:"__typename"`
	Core     struct {
		UserResults struct {
			Result struct {
				IsBlueVerified bool       `json:"is_blue_verified"`
				Legacy         legacyUser `json:"legacy"`
			} `json:"result"`
		} `json:"user_results"`
	} `json:"core"`
	Views struct {
		Count string `json:"count"`
	} `json:"views"`
	NoteTweet struct {
		NoteTweetResults struct {
			Result struct {
				Text string `json:"text"`
			} `json:"result"`
		} `json:"note_tweet_results"`
	} `json:"note_tweet"`
	QuotedStatusResult struct {
		Result *result `json:"result"`
	} `json:"quoted_status_result"`
	Legacy legacyTweet `json:"legacy"`
}

func (result *result) parse() *Tweet {
	if result.NoteTweet.NoteTweetResults.Result.Text != "" {
		result.Legacy.FullText = result.NoteTweet.NoteTweetResults.Result.Text
	}
	tw := parseLegacyTweet(&result.Core.UserResults.Result.Legacy, &result.Legacy)
	if tw.Views == 0 && result.Views.Count != "" {
		tw.Views, _ = strconv.Atoi(result.Views.Count)
	}
	if result.QuotedStatusResult.Result != nil {
		tw.QuotedStatus = result.QuotedStatusResult.Result.parse()
	}
	return tw
}

type entry struct {
	Content struct {
		CursorType string `json:"cursorType"`
		Value      string `json:"value"`
		Items      []struct {
			Item struct {
				ItemContent struct {
					TweetDisplayType string `json:"tweetDisplayType"`
					TweetResults     struct {
						Result result `json:"result"`
					} `json:"tweet_results"`
				} `json:"itemContent"`
			} `json:"item"`
		} `json:"items"`
		ItemContent struct {
			TweetDisplayType string `json:"tweetDisplayType"`
			TweetResults     struct {
				Result result `json:"result"`
			} `json:"tweet_results"`
		} `json:"itemContent"`
	} `json:"content"`
}

// timeline v2 JSON object
type timelineV2 struct {
	Data struct {
		User struct {
			Result struct {
				TimelineV2 struct {
					Timeline struct {
						Instructions []struct {
							Entries []entry `json:"entries"`
							Entry   entry   `json:"entry"`
							Type    string  `json:"type"`
						} `json:"instructions"`
					} `json:"timeline"`
				} `json:"timeline_v2"`
			} `json:"result"`
		} `json:"user"`
	} `json:"data"`
}

func (timeline *timelineV2) parseTweets() ([]*Tweet, string) {
	var cursor string
	var tweets []*Tweet
	for _, instruction := range timeline.Data.User.Result.TimelineV2.Timeline.Instructions {
		for _, entry := range instruction.Entries {
			if entry.Content.CursorType == "Bottom" {
				cursor = entry.Content.Value
				continue
			}
			if entry.Content.ItemContent.TweetResults.Result.Typename == "Tweet" {
				if tweet := entry.Content.ItemContent.TweetResults.Result.parse(); tweet != nil {
					tweets = append(tweets, tweet)
				}
			}
		}
	}
	return tweets, cursor
}

type threadedConversation struct {
	Data struct {
		ThreadedConversationWithInjectionsV2 struct {
			Instructions []struct {
				Type    string  `json:"type"`
				Entries []entry `json:"entries"`
				Entry   entry   `json:"entry"`
			} `json:"instructions"`
		} `json:"threaded_conversation_with_injections_v2"`
	} `json:"data"`
}

func (conversation *threadedConversation) parse() []*Tweet {
	var tweets []*Tweet
	for _, instruction := range conversation.Data.ThreadedConversationWithInjectionsV2.Instructions {
		for _, entry := range instruction.Entries {
			if entry.Content.ItemContent.TweetResults.Result.Typename == "Tweet" {
				if tweet := entry.Content.ItemContent.TweetResults.Result.parse(); tweet != nil {
					if entry.Content.ItemContent.TweetDisplayType == "SelfThread" {
						tweet.IsSelfThread = true
					}
					tweets = append(tweets, tweet)
				}
			}
			for _, item := range entry.Content.Items {
				if item.Item.ItemContent.TweetResults.Result.Typename == "Tweet" {
					if tweet := item.Item.ItemContent.TweetResults.Result.parse(); tweet != nil {
						if item.Item.ItemContent.TweetDisplayType == "SelfThread" {
							tweet.IsSelfThread = true
						}
						tweets = append(tweets, tweet)
					}
				}
			}
		}
	}
	for _, tweet := range tweets {
		if tweet.InReplyToStatusID != "" {
			for _, parentTweet := range tweets {
				if parentTweet.ID == tweet.InReplyToStatusID {
					tweet.InReplyToStatus = parentTweet
					break
				}
			}
		}
		if tweet.IsSelfThread && tweet.ConversationID == tweet.ID {
			for _, childTweet := range tweets {
				if childTweet.IsSelfThread && childTweet.ID != tweet.ID {
					tweet.Thread = append(tweet.Thread, childTweet)
				}
			}
			if len(tweet.Thread) == 0 {
				tweet.IsSelfThread = false
			}
		}
	}
	return tweets
}

func parseLegacyTweet(user *legacyUser, tweet *legacyTweet) *Tweet {
	username := user.ScreenName
	name := user.Name
	tweetID := tweet.IDStr
	tw := &Tweet{
		ConversationID: tweet.ConversationIDStr,
		ID:             tweetID,
		Likes:          tweet.FavoriteCount,
		Name:           name,
		PermanentURL:   fmt.Sprintf("https://twitter.com/%s/status/%s", username, tweetID),
		Replies:        tweet.ReplyCount,
		Retweets:       tweet.RetweetCount,
		Text:           tweet.FullText,
		UserID:         tweet.UserIDStr,
		Username:       username,
	}

	tm, err := time.Parse(time.RubyDate, tweet.CreatedAt)
	if err == nil {
		tw.TimeParsed = tm
		tw.Timestamp = tm.Unix()
	}

	if tweet.Place.ID != "" {
		tw.Place = &tweet.Place
	}

	if tweet.QuotedStatusIDStr != "" {
		tw.IsQuoted = true
		tw.QuotedStatusID = tweet.QuotedStatusIDStr
	}
	if tweet.InReplyToStatusIDStr != "" {
		tw.IsReply = true
		tw.InReplyToStatusID = tweet.InReplyToStatusIDStr
	}
	if tweet.RetweetedStatusIDStr != "" || tweet.RetweetedStatusResult.Result != nil {
		tw.IsRetweet = true
		tw.RetweetedStatusID = tweet.RetweetedStatusIDStr
		if tweet.RetweetedStatusResult.Result != nil {
			tw.RetweetedStatus = parseLegacyTweet(&tweet.RetweetedStatusResult.Result.Core.UserResults.Result.Legacy, &tweet.RetweetedStatusResult.Result.Legacy)
			tw.RetweetedStatusID = tw.RetweetedStatus.ID
		}
	}

	if tweet.Views.Count != "" {
		views, viewsErr := strconv.Atoi(tweet.Views.Count)
		if viewsErr != nil {
			views = 0
		}
		tw.Views = views
	}

	for _, pinned := range user.PinnedTweetIdsStr {
		if tweet.IDStr == pinned {
			tw.IsPin = true
			break
		}
	}

	for _, hash := range tweet.Entities.Hashtags {
		tw.Hashtags = append(tw.Hashtags, hash.Text)
	}

	for _, mention := range tweet.Entities.UserMentions {
		tw.Mentions = append(tw.Mentions, Mention{
			ID:       mention.IDStr,
			Username: mention.ScreenName,
			Name:     mention.Name,
		})
	}

	for _, media := range tweet.ExtendedEntities.Media {
		if media.Type == "photo" {
			photo := Photo{
				ID:  media.IDStr,
				URL: media.MediaURLHttps,
			}

			tw.Photos = append(tw.Photos, photo)
		} else if media.Type == "video" {
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
		} else if media.Type == "animated_gif" {
			gif := GIF{
				ID:      media.IDStr,
				Preview: media.MediaURLHttps,
			}

			// Twitter's API doesn't provide bitrate for GIFs, (it's always set to zero).
			// Therefore we check for `>=` instead of `>` in the loop below.
			// Also, GIFs have just a single variant today. Just in case that changes in the future,
			// and there will be multiple variants, we'll pick the one with the highest bitrate,
			// if other one will have a non-zero bitrate.
			maxBitrate := 0
			for _, variant := range media.VideoInfo.Variants {
				if variant.Bitrate >= maxBitrate {
					gif.URL = variant.URL
					maxBitrate = variant.Bitrate
				}
			}

			tw.GIFs = append(tw.GIFs, gif)
		}

		if !tw.SensitiveContent {
			sensitive := media.ExtSensitiveMediaWarning
			tw.SensitiveContent = sensitive.AdultContent || sensitive.GraphicViolence || sensitive.Other
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
	var foundedMedia []string
	tw.HTML = reTwitterURL.ReplaceAllStringFunc(tw.HTML, func(tco string) string {
		for _, entity := range tweet.Entities.URLs {
			if tco == entity.URL {
				return fmt.Sprintf(`<a href="%s">%s</a>`, entity.ExpandedURL, tco)
			}
		}
		for _, entity := range tweet.ExtendedEntities.Media {
			if tco == entity.URL {
				foundedMedia = append(foundedMedia, entity.MediaURLHttps)
				return fmt.Sprintf(`<br><a href="%s"><img src="%s"/></a>`, tco, entity.MediaURLHttps)
			}
		}
		return tco
	})
	for _, photo := range tw.Photos {
		url := photo.URL
		if stringInSlice(url, foundedMedia) {
			continue
		}
		tw.HTML += fmt.Sprintf(`<br><img src="%s"/>`, url)
	}
	for _, video := range tw.Videos {
		url := video.Preview
		if stringInSlice(url, foundedMedia) {
			continue
		}
		tw.HTML += fmt.Sprintf(`<br><img src="%s"/>`, url)
	}
	for _, gif := range tw.GIFs {
		url := gif.Preview
		if stringInSlice(url, foundedMedia) {
			continue
		}
		tw.HTML += fmt.Sprintf(`<br><img src="%s"/>`, url)
	}
	tw.HTML = strings.Replace(tw.HTML, "\n", "<br>", -1)
	return tw
}
