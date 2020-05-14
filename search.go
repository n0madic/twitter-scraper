package twitterscraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const ajaxSearchURL = "https://twitter.com/i/search/timeline?q=%s"

// GetTweets returns channel with tweets for a given search query
func GetSearchTweets(query string, maxTweetsNbr int) <-chan *Result {
	channel := make(chan *Result)
	go func(query string) {
		defer close(channel)
		var maxId string
		tweetsNbr := 0
		for tweetsNbr < maxTweetsNbr {
			tweets, err := FetchSearchTweets(query, maxId)
			if err != nil {
				channel <- &Result{Error: err}
				return
			}

			if len(tweets) == 0 {
				break
			}

			for _, tweet := range tweets {
				if tweetsNbr < maxTweetsNbr {
					lastId, _ := strconv.ParseInt(tweet.ID, 10, 64)
					maxId = strconv.FormatInt(lastId - 1, 10)
					channel <- &Result{Tweet: *tweet}
				}
				tweetsNbr++
			}
		}
	}(query)
	return channel
}

// FetchTweets gets tweets for a given search query, via the Twitter frontend API
func FetchSearchTweets(query, maxId string) ([]*Tweet, error) {
	if maxId != "" {
		query = query + " max_id:" + maxId
	}

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(ajaxSearchURL, url.PathEscape(query)),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", "https://twitter.com/search/timeline")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("X-Twitter-Active-User", "yes")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	q := req.URL.Query()
	q.Add("f", "tweets")
	q.Add("include_available_features", "1")
	q.Add("include_entities", "1")
	q.Add("include_new_items_bar", "true")

	req.URL.RawQuery = q.Encode()

	htm, err := getHTMLFromJSON(req, "items_html")
	if err != nil {
		return nil, err
	}

	tweets, err := readTweetsFromHTML(htm)
	if err != nil {
		return nil, err
	}

	return tweets, nil
}
