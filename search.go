package twitterscraper

import (
	"context"
	"errors"
	"strconv"
)

const searchURL = "https://api.twitter.com/2/search/adaptive.json"

// SearchTweets returns channel with tweets for a given search query
func (s *Scraper) SearchTweets(ctx context.Context, query string, maxTweetsNbr int) <-chan *TweetResult {
	return getTweetTimeline(ctx, query, maxTweetsNbr, s.FetchSearchTweets)
}

// SearchProfiles returns channel with profiles for a given search query
func (s *Scraper) SearchProfiles(ctx context.Context, query string, maxProfilesNbr int) <-chan *ProfileResult {
	return getUserTimeline(ctx, query, maxProfilesNbr, s.FetchSearchProfiles)
}

// getSearchTimeline gets results for a given search query, via the Twitter frontend API
func (s *Scraper) getSearchTimeline(query string, maxNbr int, cursor string) (*timelineV1, error) {
	if !s.isLogged {
		return nil, errors.New("scraper is not logged in for search")
	}

	if maxNbr > 50 {
		maxNbr = 50
	}

	req, err := s.newRequest("GET", searchURL)
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
		q.Add("tweet_search_mode", "live")
	case SearchPhotos:
		q.Add("result_filter", "image")
	case SearchVideos:
		q.Add("result_filter", "video")
	case SearchUsers:
		q.Add("result_filter", "user")
	}

	req.URL.RawQuery = q.Encode()

	var timeline timelineV1
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
