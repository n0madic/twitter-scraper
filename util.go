package twitterscraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func getHTMLFromJSON(req *http.Request, field string) (*strings.Reader, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status: %s", resp.Status)
	}

	ajaxJSON := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&ajaxJSON)
	if err != nil {
		return nil, err
	}

	htm, ok := ajaxJSON[field].(string)
	if !ok {
		return nil, fmt.Errorf("filed not found in JSON")
	}

	return strings.NewReader(htm), nil
}
