package twitterscraper

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const ajaxURL = "https://twitter.com/i/profiles/show/%s/timeline/tweets"

// Video type.
type Video struct {
	ID      string
	Preview string
}

// Tweet type.
type Tweet struct {
	Hashtags     []string
	HTML         string
	ID           string
	IsPin        bool
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
type Result struct {
	Tweet
	Error error
}

// GetTweets returns channel with tweets for a given user.
func GetTweets(ctx context.Context, user string, maxTweetsNbr int) <-chan *Result {
	channel := make(chan *Result)
	go func(user string) {
		defer close(channel)
		var lastTweetID string
		tweetsNbr := 0
		for tweetsNbr < maxTweetsNbr {
			select {
			case <-ctx.Done():
				channel <- &Result{Error: ctx.Err()}
				return
			default:
			}

			tweets, err := FetchTweets(user, lastTweetID)
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
					lastId, _ := strconv.ParseInt(tweet.ID, 10, 64)
					lastTweetID = strconv.FormatInt(lastId-1, 10)
					channel <- &Result{Tweet: *tweet}
				}
				tweetsNbr++
			}
		}
	}(user)
	return channel
}

// FetchTweets gets tweets for a given user, via the Twitter frontend API.
func FetchTweets(user string, last string) ([]*Tweet, error) {
	req, err := http.NewRequest("GET", "https://syndication.twitter.com/timeline/profile", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", "https://publish.twitter.com/")

	q := req.URL.Query()
	q.Add("screen_name", user)
	q.Add("with_replies", "true")
	if last != "" {
		q.Add("max_position", last)
	}
	req.URL.RawQuery = q.Encode()

	htm, err := getHTMLFromJSON(req, "body")
	if err != nil {
		return nil, err
	}

	tweets, err := readTweetsFromHTML(htm)
	if err != nil {
		return nil, err
	}

	return tweets, nil
}

func readTweetsFromHTML(htm *strings.Reader) ([]*Tweet, error) {
	var tweets []*Tweet

	doc, err := goquery.NewDocumentFromReader(htm)
	if err != nil {
		return nil, err
	}

	doc.Find(".timeline-Tweet").Each(func(i int, s *goquery.Selection) {
		var tweet Tweet
		timeStr, ok := s.Find(".timeline-Tweet-metadata > a > time").Attr("datetime")
		if ok {
			tweet.TimeParsed, _ = time.Parse("2006-01-02T15:04:05-0700", timeStr)
			tweet.Timestamp = tweet.TimeParsed.Unix()
			tweet.ID = s.AttrOr("data-tweet-id", "")
			//tweet.UserID = s.Find(".tweet").AttrOr("data-user-id", "")
			tweet.Username = strings.TrimPrefix(s.Find(".TweetAuthor-screenName").AttrOr("title", ""), "@")
			tweet.PermanentURL = fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.Username, tweet.ID)
			tweet.Text = s.Find(".timeline-Tweet-text").Text()
			tweet.HTML, _ = s.Find(".timeline-Tweet-text").Html()
			s.Find(".timeline-Tweet-retweetCredit").Each(func(i int, c *goquery.Selection) {
				tweet.IsRetweet = true
			})
			// s.Find("span.js-pinned-text").Each(func(i int, c *goquery.Selection) {
			// 	tweet.IsPin = true
			// })
			// s.Find(".ProfileTweet-actionCount").Each(func(i int, c *goquery.Selection) {
			// 	txt := strings.TrimSpace(c.Text())
			// 	switch {
			// 	case strings.HasSuffix(txt, "likes"):
			// 		l := strings.Split(txt, " ")
			// 		tweet.Likes, _ = strconv.Atoi(l[0])
			// 	case strings.HasSuffix(txt, "replies"):
			// 		l := strings.Split(txt, " ")
			// 		tweet.Replies, _ = strconv.Atoi(l[0])
			// 	case strings.HasSuffix(txt, "retweets"):
			// 		l := strings.Split(txt, " ")
			// 		tweet.Retweets, _ = strconv.Atoi(l[0])
			// 	}
			// })
			s.Find(".hashtag > span.PrettyLink-value").Each(func(i int, h *goquery.Selection) {
				tweet.Hashtags = append(tweet.Hashtags, h.Text())
			})
			s.Find("a.link:not(.u-hidden)").Each(func(i int, u *goquery.Selection) {
				if link, ok := u.Attr("data-expanded-url"); ok {
					tweet.URLs = append(tweet.URLs, link)
				}
			})
			s.Find(".NaturalImage-image").Each(func(i int, p *goquery.Selection) {
				if link, ok := p.Attr("data-image"); ok {
					tweet.Photos = append(tweet.Photos, link+"?format=jpg&name=large")
				}
			})
			// s.Find(".PlayableMedia-player").Each(func(i int, v *goquery.Selection) {
			// 	if style, ok := v.Attr("style"); ok {
			// 		if strings.Contains(style, "background") {
			// 			match := regexp.MustCompile(`https:\/\/.+\/([\w-]+)\.(?:jpg|png)`).FindStringSubmatch(style)
			// 			if len(match) == 2 {
			// 				tweet.Videos = append(tweet.Videos, Video{ID: match[1], Preview: match[0]})
			// 			}
			// 		}
			// 	}
			// })
			tweets = append(tweets, &tweet)
		}
	})

	return tweets, nil
}
