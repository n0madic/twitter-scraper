package twitterscraper

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	reHashtag    = regexp.MustCompile(`\B(\#\S+\b)`)
	reTwitterURL = regexp.MustCompile(`https:(\/\/t\.co\/([A-Za-z0-9]|[A-Za-z]){10})`)
	reUsername   = regexp.MustCompile(`\B(\@\S{1,15}\b)`)
)

func (s *Scraper) newRequest(method string, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("include_profile_interstitial_type", "1")
	q.Add("include_blocking", "1")
	q.Add("include_blocked_by", "1")
	q.Add("include_followed_by", "1")
	q.Add("include_want_retweets", "1")
	q.Add("include_mute_edge", "1")
	q.Add("include_can_dm", "1")
	q.Add("include_can_media_tag", "1")
	q.Add("skip_status", "1")
	q.Add("cards_platform", "Web-12")
	q.Add("include_cards", "1")
	q.Add("include_ext_alt_text", "true")
	q.Add("include_quote_count", "true")
	q.Add("include_reply_count", "1")
	q.Add("tweet_mode", "extended")
	q.Add("include_entities", "true")
	q.Add("include_user_entities", "true")
	q.Add("include_ext_media_color", "true")
	q.Add("include_ext_media_availability", "true")
	q.Add("send_error_codes", "true")
	q.Add("simple_quoted_tweet", "true")
	q.Add("include_tweet_replies", strconv.FormatBool(s.includeReplies))
	q.Add("ext", "mediaStats,highlightedLabel")
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func getUserTimeline(ctx context.Context, query string, maxProfilesNbr int, fetchFunc fetchProfileFunc) <-chan *ProfileResult {
	channel := make(chan *ProfileResult)
	go func(query string) {
		defer close(channel)
		var nextCursor string
		profilesNbr := 0
		for profilesNbr < maxProfilesNbr {
			select {
			case <-ctx.Done():
				channel <- &ProfileResult{Error: ctx.Err()}
				return
			default:
			}

			profiles, next, err := fetchFunc(query, maxProfilesNbr, nextCursor)
			if err != nil {
				channel <- &ProfileResult{Error: err}
				return
			}

			if len(profiles) == 0 {
				break
			}

			for _, profile := range profiles {
				select {
				case <-ctx.Done():
					channel <- &ProfileResult{Error: ctx.Err()}
					return
				default:
				}

				if profilesNbr < maxProfilesNbr {
					nextCursor = next
					channel <- &ProfileResult{Profile: *profile}
				} else {
					break
				}
				profilesNbr++
			}
		}
	}(query)
	return channel
}

func getTweetTimeline(ctx context.Context, query string, maxTweetsNbr int, fetchFunc fetchTweetFunc) <-chan *TweetResult {
	channel := make(chan *TweetResult)
	go func(query string) {
		defer close(channel)
		var nextCursor string
		tweetsNbr := 0
		for tweetsNbr < maxTweetsNbr {
			select {
			case <-ctx.Done():
				channel <- &TweetResult{Error: ctx.Err()}
				return
			default:
			}

			tweets, next, err := fetchFunc(query, maxTweetsNbr, nextCursor)
			if err != nil {
				channel <- &TweetResult{Error: err}
				return
			}

			if len(tweets) == 0 {
				break
			}

			for _, tweet := range tweets {
				select {
				case <-ctx.Done():
					channel <- &TweetResult{Error: ctx.Err()}
					return
				default:
				}

				if tweetsNbr < maxTweetsNbr {
					if tweet.IsPin && nextCursor != "" {
						continue
					}
					nextCursor = next
					channel <- &TweetResult{Tweet: *tweet}
				} else {
					break
				}
				tweetsNbr++
			}
		}
	}(query)
	return channel
}

func parseProfile(user legacyUser) Profile {
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

	return profile
}

func parseTimeline(timeline *timeline) ([]*Tweet, string) {
	tweets := make(map[string]Tweet)

	for id, tweet := range timeline.GlobalObjects.Tweets {
		username := timeline.GlobalObjects.Users[tweet.UserIDStr].ScreenName
		tw := Tweet{
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
		}
		if tweet.InReplyToStatusIDStr != "" {
			tw.IsReply = true
		}
		if tweet.RetweetedStatusIDStr != "" {
			tw.IsRetweet = true
			if retweet, ok := timeline.GlobalObjects.Tweets[tweet.RetweetedStatusIDStr]; ok {
				tw.Retweet = Retweet{
					ID:       tweet.RetweetedStatusIDStr,
					UserID:   retweet.UserIDStr,
					Username: timeline.GlobalObjects.Users[retweet.UserIDStr].ScreenName,
				}
				tm, err := time.Parse(time.RubyDate, retweet.CreatedAt)
				if err == nil {
					tw.Retweet.TimeParsed = tm
					tw.Retweet.Timestamp = tm.Unix()
				}
			}
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

		tweets[tw.ID] = tw
	}

	var cursor string
	var pinnedTweet *Tweet
	var orderedTweets []*Tweet
	for _, instruction := range timeline.Timeline.Instructions {
		if instruction.PinEntry.Entry.Content.Item.Content.Tweet.ID != "" {
			if tweet, ok := tweets[instruction.PinEntry.Entry.Content.Item.Content.Tweet.ID]; ok {
				pinnedTweet = &tweet
			}
		}
		for _, entry := range instruction.AddEntries.Entries {
			if tweet, ok := tweets[entry.Content.Item.Content.Tweet.ID]; ok {
				orderedTweets = append(orderedTweets, &tweet)
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

func parseUsers(timeline *timeline) ([]*Profile, string) {
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
