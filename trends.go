package twitterscraper

import "fmt"

// GetTrends return list of trends.
func (s *Scraper) GetTrends() ([]string, error) {
	req, err := s.newRequest("GET", "https://api.twitter.com/2/guide.json")
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("count", "20")
	q.Add("candidate_source", "trends")
	q.Add("include_page_configuration", "false")
	q.Add("entity_tokens", "false")
	req.URL.RawQuery = q.Encode()

	var jsn timeline
	curBearerToken := s.bearerToken
	if curBearerToken != bearerToken2 {
		s.setBearerToken(bearerToken2)
	}
	err = s.RequestAPI(req, &jsn)
	if curBearerToken != bearerToken2 {
		s.setBearerToken(curBearerToken)
	}
	if err != nil {
		return nil, err
	}

	if len(jsn.Timeline.Instructions[1].AddEntries.Entries) < 2 {
		return nil, fmt.Errorf("no trend entries found")
	}

	var trends []string
	for _, item := range jsn.Timeline.Instructions[1].AddEntries.Entries[1].Content.TimelineModule.Items {
		trends = append(trends, item.Item.ClientEventInfo.Details.GuideDetails.TransparentGuideDetails.TrendMetadata.TrendName)
	}

	return trends, nil
}
