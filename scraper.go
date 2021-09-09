package twitterscraper

import (
	"crypto/tls"
	"errors"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Scraper object
type Scraper struct {
	client         *http.Client
	clientTimeout  time.Duration
	delay          int64
	guestToken     string
	guestCreatedAt time.Time
	includeReplies bool
	searchMode     SearchMode
	wg             sync.WaitGroup

	cookie     string
	xCsrfToken string
}

// SearchMode type
type SearchMode int

const (
	// SearchTop - default mode
	SearchTop SearchMode = iota
	// SearchLatest - live mode
	SearchLatest
	// SearchPhotos - image mode
	SearchPhotos
	// SearchVideos - video mode
	SearchVideos
	// SearchUsers - user mode
	SearchUsers
)

// default http client timeout
const DefaultClientTimeout = 10 * time.Second

var defaultScraper *Scraper

// New creates a Scraper object
func New() *Scraper {
	return &Scraper{
		client:        &http.Client{Timeout: DefaultClientTimeout},
		clientTimeout: DefaultClientTimeout,
	}
}

// SetSearchMode switcher
func (s *Scraper) SetSearchMode(mode SearchMode) *Scraper {
	s.searchMode = mode
	return s
}

// SetSearchMode wrapper for default Scraper
func SetSearchMode(mode SearchMode) *Scraper {
	return defaultScraper.SetSearchMode(mode)
}

// WithDelay add delay between API requests (in seconds)
func (s *Scraper) WithDelay(seconds int64) *Scraper {
	s.delay = seconds
	return s
}

// WithDelay wrapper for default Scraper
func WithDelay(seconds int64) *Scraper {
	return defaultScraper.WithDelay(seconds)
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

// cookie
func (s *Scraper) WithCookie(cookie string) *Scraper {
	s.cookie = cookie
	return s
}

// x csrf token
func (s *Scraper) WithXCsrfToken(xcsrfToken string) *Scraper {
	s.xCsrfToken = xcsrfToken
	return s
}

// client timeout
func (s *Scraper) WithClientTimeout(timeout time.Duration) *Scraper {
	s.clientTimeout = timeout
	return s
}

// SetProxy
// set http proxy in the format `http://HOST:PORT`
// set socket proxy in the format `socks5://HOST:PORT`
func (s *Scraper) SetProxy(proxyAddr string) error {
	if strings.HasPrefix(proxyAddr, "http") {
		urlproxy, err := url.Parse(proxyAddr)
		if err != nil {
			return err
		}
		s.client = &http.Client{
			Transport: &http.Transport{
				Proxy:        http.ProxyURL(urlproxy),
				TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
				DialContext: (&net.Dialer{
					Timeout: s.clientTimeout,
				}).DialContext,
			},
		}
		return nil
	}
	if strings.HasPrefix(proxyAddr, "socks5") {
		baseDialer := &net.Dialer{
			Timeout:   s.clientTimeout,
			KeepAlive: s.clientTimeout,
		}
		socksHostPort := strings.ReplaceAll(proxyAddr, "socks5://", "")
		dialSocksProxy, err := proxy.SOCKS5("tcp", socksHostPort, nil, baseDialer)
		if err != nil {
			return errors.New("error creating socks5 proxy :" + err.Error())
		}
		if contextDialer, ok := dialSocksProxy.(proxy.ContextDialer); ok {
			dialContext := contextDialer.DialContext
			s.client = &http.Client{
				Transport: &http.Transport{
					DialContext: dialContext,
				},
			}
		} else {
			return errors.New("failed type assertion to DialContext")
		}
		return nil
	}
	return errors.New("only support http(s) or socks5 protocol")
}

// SetProxy wrapper for default Scraper
func SetProxy(proxy string) error {
	return defaultScraper.SetProxy(proxy)
}

func init() {
	defaultScraper = New()
}
