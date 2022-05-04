package twitterscraper

import (
	"context"
	"fmt"
	"strconv"
)

// GetTweets returns channel with tweets for a given user.
func (s *Scraper) GetTweets(ctx context.Context, user string, maxTweetsNbr int) <-chan *TweetResult {
	return getTweetTimeline(ctx, user, maxTweetsNbr, s.FetchTweets)
}

// Deprecated: GetTweets wrapper for default Scraper
func GetTweets(ctx context.Context, user string, maxTweetsNbr int) <-chan *TweetResult {
	return defaultScraper.GetTweets(ctx, user, maxTweetsNbr)
}

// FetchTweets gets tweets for a given user, via the Twitter frontend API.
func (s *Scraper) FetchTweets(user string, maxTweetsNbr int, cursor string) ([]*Tweet, string, error) {
	if maxTweetsNbr > 200 {
		maxTweetsNbr = 200
	}

	userID, err := s.GetUserIDByScreenName(user)
	if err != nil {
		return nil, "", err
	}

	req, err := s.newRequest("GET", "https://api.twitter.com/2/timeline/profile/"+userID+".json")
	if err != nil {
		return nil, "", err
	}

	q := req.URL.Query()
	q.Add("count", strconv.Itoa(maxTweetsNbr))
	q.Add("userId", userID)
	if cursor != "" {
		q.Add("cursor", cursor)
	}
	req.URL.RawQuery = q.Encode()

	var timeline timeline
	err = s.RequestAPI(req, &timeline)
	if err != nil {
		return nil, "", err
	}

	tweets, nextCursor := timeline.parseTweets()
	return tweets, nextCursor, nil
}

// GetTweet get a single tweet by ID.
func (s *Scraper) GetTweet(id string) (*Tweet, error) {
	req, err := s.newRequest("GET", "https://twitter.com/i/api/2/timeline/conversation/"+id+".json")
	if err != nil {
		return nil, err
	}

	var timeline timeline
	err = s.RequestAPI(req, &timeline)
	if err != nil {
		return nil, err
	}

	tweets, _ := timeline.parseTweets()
	for _, tweet := range tweets {
		if tweet.ID == id {
			return tweet, nil
		}
	}
	return nil, fmt.Errorf("tweet with ID %s not found", id)
}

// Deprecated: GetTweet wrapper for default Scraper
func GetTweet(id string) (*Tweet, error) {
	return defaultScraper.GetTweet(id)
}
