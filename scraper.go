package twitterscraper

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Scraper object
type Scraper struct {
	client         *http.Client
	guestToken     string
	includeReplies bool
	searchMode     string
}

var defaultScraper *Scraper

// New creates a Scraper object
func New() *Scraper {
	return &Scraper{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// SetSearchLive enable/disable realtime search
func (s *Scraper) SetSearchLive(srctype bool) *Scraper {
	if srctype {
		s.searchMode = "live"
	}
	return s
}

// SetSearchLive wrapper for default SetSearchLive
func SetSearchLive(srctype bool) *Scraper {
	return defaultScraper.SetSearchLive(srctype)
}

// SetSearchPhotos filter search for photos only
func (s *Scraper) SetSearchPhotos(srctype bool) *Scraper {
	if srctype {
		s.searchMode = "image"
	}
	return s
}

// SetSearchPhotos wrapper for default SetSearchPhotos
func SetSearchPhotos(srctype bool) *Scraper {
	return defaultScraper.SetSearchPhotos(srctype)
}

// SetSearchVideos filter search for videos only
func (s *Scraper) SetSearchVideos(srctype bool) *Scraper {
	if srctype {
		s.searchMode = "video"
	}
	return s
}

// SetSearchVideos wrapper for default SetSearchVideos
func SetSearchVideos(srctype bool) *Scraper {
	return defaultScraper.SetSearchVideos(srctype)
}

// WithReplies enable/disable load timeline with tweet replies
func (s *Scraper) WithReplies(b bool) *Scraper {
	s.includeReplies = b
	return s
}

// WithReplies wrapper for default Scraper
func WithReplies(b bool) *Scraper {
	return defaultScraper.WithReplies(b)
}

// SetProxy set http proxy in the format `http://HOST:PORT`
func (s *Scraper) SetProxy(proxy string) error {
	if !strings.HasPrefix(proxy, "http") {
		return errors.New("only support http(s) protocol")
	}
	urlproxy, err := url.Parse(proxy)
	if err != nil {
		return err
	}
	s.client = &http.Client{
		Transport: &http.Transport{
			Proxy:        http.ProxyURL(urlproxy),
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			DialContext: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).DialContext,
		},
	}
	return nil
}

// SetProxy wrapper for default Scraper
func SetProxy(proxy string) error {
	return defaultScraper.SetProxy(proxy)
}

func init() {
	defaultScraper = New()
}
