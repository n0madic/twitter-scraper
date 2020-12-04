package twitterscraper

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const mobileSearchURL = "https://mobile.twitter.com/search?q=%s"

// SearchTweets returns channel with tweets for a given search query
func SearchTweets(ctx context.Context, query string, maxTweetsNbr int) <-chan *Result {
	channel := make(chan *Result)
	go func(query string) {
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

			tweets, next, err := FetchSearchTweets(query, nextCursor)
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

// FetchSearchTweets gets tweets for a given search query, via the Twitter frontend API
func FetchSearchTweets(query, nextCursor string) ([]*Tweet, string, error) {
	url := fmt.Sprintf(mobileSearchURL, url.PathEscape(query))
	if nextCursor != "" {
		url = "https://mobile.twitter.com" + nextCursor
	}

	client := http.DefaultClient
	if HTTPProxy != nil {
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(HTTPProxy),
				DialContext: (&net.Dialer{
					Timeout: 10 * time.Second,
				}).DialContext,
			},
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	req.Header.Set("Referer", "https://mobile.twitter.com/")
	req.Header.Set("User-Agent", "Opera/9.80 (J2ME/MIDP; Opera Mini/5.1.21214/28.2725; U; ru) Presto/2.8.119 Version/11.10")

	resp, err := client.Do(req)
	if resp == nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("response status: %s", resp.Status)
	}

	return readTweetsFromMobileHTML(resp.Body)
}

func readTweetsFromMobileHTML(htm io.ReadCloser) ([]*Tweet, string, error) {
	var tweets []*Tweet

	doc, err := goquery.NewDocumentFromReader(htm)
	if err != nil {
		return nil, "", err
	}

	doc.Find("table.tweet").Each(func(i int, s *goquery.Selection) {
		var tweet Tweet
		tweetID, ok := s.Find(".tweet-text").Attr("data-id")
		if ok {
			tweet.ID = tweetID
			tweet.Username = strings.TrimPrefix(strings.TrimSpace(s.Find("td.user-info > a > div.username").Text()), "@")
			tweet.PermanentURL = fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.Username, tweet.ID)
			tweet.Text = strings.TrimSpace(s.Find(".tweet-text").Text())
			tweet.HTML, _ = s.Find(".tweet-text").Html()
			tweet.HTML = strings.TrimSpace(tweet.HTML)
			s.Find("td.tweet-social-context > span").Each(func(i int, c *goquery.Selection) {
				tweet.IsRetweet = true
			})
			s.Find(".twitter-hashtag").Each(func(i int, h *goquery.Selection) {
				tweet.Hashtags = append(tweet.Hashtags, h.Text())
			})
			s.Find("a.tco-link:not(.u-hidden)").Each(func(i int, u *goquery.Selection) {
				if link, ok := u.Attr("data-expanded-url"); ok {
					tweet.URLs = append(tweet.URLs, link)
				}
			})
			s.Find("div.media > img").Each(func(i int, p *goquery.Selection) {
				if link, ok := p.Attr("src"); ok {
					tweet.Photos = append(tweet.Photos, strings.TrimSuffix(link, ":small"))
				}
			})

			tweets = append(tweets, &tweet)
		}
	})

	nextCursor := doc.Find("div.w-button-more > a").AttrOr("href", "")

	return tweets, nextCursor, nil
}
