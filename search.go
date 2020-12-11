package twitterscraper

import (
	"context"
	"net/url"
	"strconv"
)

// SearchTweets returns channel with tweets for a given search query
func SearchTweets(ctx context.Context, query string, maxTweetsNbr int) <-chan *Result {
	return getTimeline(ctx, query, maxTweetsNbr, FetchSearchTweets)
}

// FetchSearchTweets gets tweets for a given search query, via the Twitter frontend API
func FetchSearchTweets(query string, maxTweetsNbr int, cursor string) ([]*Tweet, string, error) {
	query = url.PathEscape(query)
	if maxTweetsNbr > 200 {
		maxTweetsNbr = 200
	}

	req, err := newRequest("GET", "https://twitter.com/i/api/2/search/adaptive.json")
	if err != nil {
		return nil, "", err
	}

	q := req.URL.Query()
	q.Add("q", query)
	q.Add("count", strconv.Itoa(maxTweetsNbr))
	q.Add("query_source", "typed_query")
	q.Add("pc", "1")
	q.Add("spelling_corrections", "1")
	if cursor != "" {
		q.Add("cursor", cursor)
	}
	req.URL.RawQuery = q.Encode()

	var timeline timeline
	err = requestAPI(req, &timeline)
	if err != nil {
		return nil, "", err
	}

	tweets, nextCursor := parseTimeline(&timeline)
	return tweets, nextCursor, nil
}
