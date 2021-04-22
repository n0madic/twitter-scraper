package twitterscraper

import (
	"testing"
)

func TestGetTrends(t *testing.T) {
	trends, err := GetTrends()
	if err != nil {
		t.Error(err)
	}

	if len(trends) != 20 {
		t.Errorf("Expected 20 trends, got %d: %#v", len(trends), trends)
	}

	for _, trend := range trends {
		if trend == "" {
			t.Error("Expected trend is empty")
		}
	}
}
