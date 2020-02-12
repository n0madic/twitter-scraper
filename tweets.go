package twitterscraper

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const ajaxURL = "https://twitter.com/i/profiles/show/%s/timeline/tweets"

// Video type
type Video struct {
	ID      string
	Preview string
}

// Tweet type
type Tweet struct {
	Hashtags     []string
	HTML         string
	ID           string
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
	Videos       []Video
}

// Result of scrapping
type Result struct {
	Tweet
	Error error
}

// GetTweets returns channel with tweets for a given user
func GetTweets(user string, pages int) <-chan *Result {
	channel := make(chan *Result, 0)
	go func(user string) {
		defer close(channel)
		var lastTweetID string
		for pages > 0 {
			tweets, err := FetchTweets(user, lastTweetID)
			if err != nil {
				channel <- &Result{Error: err}
				return
			}
			for _, tweet := range tweets {
				lastTweetID = tweet.ID
				channel <- &Result{Tweet: *tweet}
			}
			pages--
		}
	}(user)
	return channel
}

// FetchTweets gets tweets for a given user, via the Twitter frontend API
func FetchTweets(user string, last string) ([]*Tweet, error) {
	var tweets []*Tweet

	req, err := http.NewRequest("GET", fmt.Sprintf(ajaxURL, user), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", "https://twitter.com/"+user)
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8")
	req.Header.Set("X-Twitter-Active-User", "yes")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	q := req.URL.Query()
	q.Add("include_available_features", "1")
	q.Add("include_entities", "1")
	q.Add("include_new_items_bar", "true")
	if last != "" {
		q.Add("max_position", last)
	}
	req.URL.RawQuery = q.Encode()

	htm, err := getHTMLFromJSON(req, "items_html")
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(htm)
	if err != nil {
		return nil, err
	}

	doc.Find(".stream-item").Each(func(i int, s *goquery.Selection) {
		var tweet Tweet
		timeStr, ok := s.Find("._timestamp").Attr("data-time")
		if ok {
			tweet.Timestamp, _ = strconv.ParseInt(timeStr, 10, 64)
			tweet.TimeParsed = time.Unix(tweet.Timestamp, 0)
			tweet.ID = s.AttrOr("data-item-id", "")
			tweet.PermanentURL = fmt.Sprintf("https://twitter.com/%s/status/%s", user, tweet.ID)
			tweet.Text = s.Find(".tweet-text").Text()
			tweet.HTML, _ = s.Find(".tweet-text").Html()
			s.Find(".js-retweet-text, .QuoteTweet").Each(func(i int, c *goquery.Selection) {
				tweet.IsRetweet = true
			})
			s.Find(".ProfileTweet-actionCount").Each(func(i int, c *goquery.Selection) {
				txt := strings.TrimSpace(c.Text())
				if strings.HasSuffix(txt, "likes") {
					l := strings.Split(txt, " ")
					tweet.Likes, _ = strconv.Atoi(l[0])
				} else if strings.HasSuffix(txt, "replies") {
					l := strings.Split(txt, " ")
					tweet.Replies, _ = strconv.Atoi(l[0])
				} else if strings.HasSuffix(txt, "retweets") {
					l := strings.Split(txt, " ")
					tweet.Retweets, _ = strconv.Atoi(l[0])
				}
			})
			s.Find(".twitter-hashtag").Each(func(i int, h *goquery.Selection) {
				tweet.Hashtags = append(tweet.Hashtags, h.Text())
			})
			s.Find("a.twitter-timeline-link:not(.u-hidden)").Each(func(i int, u *goquery.Selection) {
				if link, ok := u.Attr("data-expanded-url"); ok {
					tweet.URLs = append(tweet.URLs, link)
				}
			})
			s.Find(".AdaptiveMedia-photoContainer").Each(func(i int, p *goquery.Selection) {
				if link, ok := p.Attr("data-image-url"); ok {
					tweet.Photos = append(tweet.Photos, link)
				}
			})
			s.Find(".PlayableMedia-player").Each(func(i int, v *goquery.Selection) {
				if style, ok := v.Attr("style"); ok {
					if strings.Contains(style, "background") {
						match := regexp.MustCompile(`https:\/\/.+\/([\w-]+)\.(?:jpg|png)`).FindStringSubmatch(style)
						if len(match) == 2 {
							tweet.Videos = append(tweet.Videos, Video{ID: match[1], Preview: match[0]})
						}
					}
				}
			})
			tweets = append(tweets, &tweet)
		}
	})

	return tweets, nil
}
