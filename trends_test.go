package twitterscraper

import (
	"testing"
)

func TestGetTrends(t *testing.T) {
	trends, err := GetTrends()
	if err != nil {
		t.Error(err)
	}

	if len(trends) != 10 {
		t.Errorf("Expected 10 trends, got %d: %#v", len(trends), trends)
	}
}
