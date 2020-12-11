package twitterscraper

// GetTrends return list of trends.
func GetTrends() ([]string, error) {
	req, err := newRequest("GET", "https://twitter.com/i/api/2/guide.json")
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
	err = requestAPI(req, &jsn)
	if err != nil {
		return nil, err
	}

	var trends []string
	for _, item := range jsn.Timeline.Instructions[1].AddEntries.Entries[1].Content.TimelineModule.Items {
		trends = append(trends, item.Item.ClientEventInfo.Details.GuideDetails.TransparentGuideDetails.TrendMetadata.TrendName)
	}

	return trends, nil
}
