package twitterscraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

//HttpProxy Public variable for Http proxy
var HTTPProxy *url.URL

//SetProxy set http proxy format `http://HOST:PORT`
func SetProxy(Proxy string) error {
	match, _ := regexp.MatchString("http.+", Proxy)
	if !match {
		return errors.New("only support http protocol")
	}
	urlproxy, err := url.Parse(Proxy)
	if err != nil {
		return err
	}
	HTTPProxy = urlproxy
	return nil
}

func newRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8")
	req.Header.Set("X-Twitter-Active-User", "yes")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	return req, nil
}

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
		return nil, fmt.Errorf("field %s not found in JSON", field)
	}

	return strings.NewReader(htm), nil
}
