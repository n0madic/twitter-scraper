package twitterscraper

import (
	"context"
	"strconv"
)

// GetTweets returns channel with tweets for a given user.
func GetTweets(ctx context.Context, user string, maxTweetsNbr int) <-chan *Result {
	return getTimeline(ctx, user, maxTweetsNbr, FetchTweets)
}

// FetchTweets gets tweets for a given user, via the Twitter frontend API.
func FetchTweets(user string, maxTweetsNbr int, cursor string) ([]*Tweet, string, error) {
	if maxTweetsNbr > 200 {
		maxTweetsNbr = 200
	}

	userID, err := GetUserIDByScreenName(user)
	if err != nil {
		return nil, "", err
	}

	req, err := newRequest("GET", "https://api.twitter.com/2/timeline/profile/"+userID+".json")
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
	err = requestAPI(req, &timeline)
	if err != nil {
		return nil, "", err
	}

	tweets, nextCursor := parseTimeline(&timeline)
	return tweets, nextCursor, nil
}
