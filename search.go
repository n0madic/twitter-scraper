package twitterscraper

import (
	"context"
	"net/url"
	"strconv"
)

// SearchTweets returns channel with tweets for a given search query
func (s *Scraper) SearchTweets(ctx context.Context, query string, maxTweetsNbr int) <-chan *Result {
	return getTimeline(ctx, query, maxTweetsNbr, s.FetchSearchTweets)
}

// SearchTweets wrapper for default Scraper
func SearchTweets(ctx context.Context, query string, maxTweetsNbr int) <-chan *Result {
	return defaultScraper.SearchTweets(ctx, query, maxTweetsNbr)
}

// FetchSearchTweets gets tweets for a given search query, via the Twitter frontend API
func (s *Scraper) FetchSearchTweets(query string, maxTweetsNbr int, cursor string) ([]*Tweet, string, error) {
	query = url.PathEscape(query)
	if maxTweetsNbr > 200 {
		maxTweetsNbr = 200
	}

	req, err := s.newRequest("GET", "https://twitter.com/i/api/2/search/adaptive.json")
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
	switch s.searchMode {
	case SearchLatest:
		q.Add("tweet_search_mode", "live")
	case SearchPhotos:
		q.Add("result_filter", "image")
	case SearchVideos:
		q.Add("result_filter", "video")
	}

	req.URL.RawQuery = q.Encode()

	var timeline timeline
	err = s.RequestAPI(req, &timeline)
	if err != nil {
		return nil, "", err
	}

	tweets, nextCursor := parseTimeline(&timeline)
	return tweets, nextCursor, nil
}
