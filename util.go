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
	reHashtag    = regexp.MustCompile(`\B(\#\w+\b)`)
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

func getTimeline(ctx context.Context, query string, maxTweetsNbr int, fetchFunc fetchFunc) <-chan *Result {
	channel := make(chan *Result)
	go func(user string) {
		defer close(channel)
		var nextCursor string
		tweetsNbr := 0
		for tweetsNbr < maxTweetsNbr {
			select {
			case <-ctx.Done():
				channel <- &Result{Error: ctx.Err()}
				return
			default:
			}

			tweets, next, err := fetchFunc(query, maxTweetsNbr, nextCursor)
			if err != nil {
				channel <- &Result{Error: err}
				return
			}

			if len(tweets) == 0 {
				break
			}

			for _, tweet := range tweets {
				select {
				case <-ctx.Done():
					channel <- &Result{Error: ctx.Err()}
					return
				default:
				}

				if tweetsNbr < maxTweetsNbr {
					nextCursor = next
					channel <- &Result{Tweet: *tweet}
				}
				tweetsNbr++
			}
		}
	}(query)
	return channel
}

func parseTimeline(timeline *timeline) ([]*Tweet, string) {
	tweets := make(map[string]Tweet)

	for id, tweet := range timeline.GlobalObjects.Tweets {
		username := timeline.GlobalObjects.Users[tweet.UserIDStr].ScreenName
		tw := Tweet{
			ID:           id,
			Likes:        tweet.FavoriteCount,
			PermanentURL: fmt.Sprintf("https://twitter.com/%s/status/%s", username, id),
			Replies:      tweet.RetweetCount,
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
	var orderedTweets []*Tweet
	if len(timeline.Timeline.Instructions) > 0 {
		for _, entry := range timeline.Timeline.Instructions[0].AddEntries.Entries {
			if tweet, ok := tweets[entry.Content.Item.Content.Tweet.ID]; ok {
				orderedTweets = append(orderedTweets, &tweet)
			}
			if entry.Content.Operation.Cursor.CursorType == "Bottom" {
				cursor = entry.Content.Operation.Cursor.Value
			}
		}
	}
	return orderedTweets, cursor
}
