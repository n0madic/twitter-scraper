package twitterscraper

import (
	"context"
	"errors"
	"strconv"
)

// SearchTweets returns channel with tweets for a given search query
func (s *Scraper) SearchTweets(ctx context.Context, query string, maxTweetsNbr int) <-chan *TweetResult {
	return getTweetTimeline(ctx, query, maxTweetsNbr, s.FetchSearchTweets)
}

// Deprecated: SearchTweets wrapper for default Scraper
func SearchTweets(ctx context.Context, query string, maxTweetsNbr int) <-chan *TweetResult {
	return defaultScraper.SearchTweets(ctx, query, maxTweetsNbr)
}

// SearchProfiles returns channel with profiles for a given search query
func (s *Scraper) SearchProfiles(ctx context.Context, query string, maxProfilesNbr int) <-chan *ProfileResult {
	return getUserTimeline(ctx, query, maxProfilesNbr, s.FetchSearchProfiles)
}

// Deprecated: SearchProfiles wrapper for default Scraper
func SearchProfiles(ctx context.Context, query string, maxProfilesNbr int) <-chan *ProfileResult {
	return defaultScraper.SearchProfiles(ctx, query, maxProfilesNbr)
}

// getSearchTimeline gets results for a given search query, via the Twitter frontend API
func (s *Scraper) getSearchTimeline(query string, maxNbr int, cursor string) (*timeline, error) {
	if !s.isLogged {
		return nil, errors.New("scraper is not logged in for search")
	}

	if maxNbr > 50 {
		maxNbr = 50
	}

	req, err := s.newRequest("GET", "https://twitter.com/i/api/2/search/adaptive.json")
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("q", query)
	q.Add("count", strconv.Itoa(maxNbr))
	q.Add("query_source", "typed_query")
	q.Add("pc", "1")
	q.Add("requestContext", "launch")
	q.Add("spelling_corrections", "1")
	q.Add("include_ext_edit_control", "true")
	if cursor != "" {
		q.Add("cursor", cursor)
	}
	switch s.searchMode {
	case SearchLatest:
		q.Add("f", "live")
	case SearchPhotos:
		q.Add("result_filter", "image")
	case SearchVideos:
		q.Add("result_filter", "video")
	case SearchUsers:
		q.Add("result_filter", "user")
	}

	req.URL.RawQuery = q.Encode()

	var timeline timeline
	err = s.RequestAPI(req, &timeline)
	if err != nil {
		return nil, err
	}
	return &timeline, nil
}

// FetchSearchTweets gets tweets for a given search query, via the Twitter frontend API
func (s *Scraper) FetchSearchTweets(query string, maxTweetsNbr int, cursor string) ([]*Tweet, string, error) {
	timeline, err := s.getSearchTimeline(query, maxTweetsNbr, cursor)
	if err != nil {
		return nil, "", err
	}
	tweets, nextCursor := timeline.parseTweets()
	return tweets, nextCursor, nil
}

// FetchSearchProfiles gets users for a given search query, via the Twitter frontend API
func (s *Scraper) FetchSearchProfiles(query string, maxProfilesNbr int, cursor string) ([]*Profile, string, error) {
	timeline, err := s.getSearchTimeline(query, maxProfilesNbr, cursor)
	if err != nil {
		return nil, "", err
	}
	users, nextCursor := timeline.parseUsers()
	return users, nextCursor, nil
}
