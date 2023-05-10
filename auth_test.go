package twitterscraper_test

import (
	"os"
	"testing"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

var (
	username = os.Getenv("TWITTER_USERNAME")
	password = os.Getenv("TWITTER_PASSWORD")
	email    = os.Getenv("TWITTER_EMAIL")
)

func TestAuth(t *testing.T) {
	scraper := twitterscraper.New()
	if err := scraper.Login(username, password, email); err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if !scraper.IsLoggedIn() {
		t.Fatalf("Expected IsLoggedIn() = true")
	}
	cookies := scraper.GetCookies()
	scraper.Logout()
	if scraper.IsLoggedIn() {
		t.Error("Expected IsLoggedIn() = false")
	}
	scraper2 := twitterscraper.New()
	scraper2.SetCookies(cookies)
	if !scraper2.IsLoggedIn() {
		t.Error("Expected restored IsLoggedIn() = true")
	}
}
