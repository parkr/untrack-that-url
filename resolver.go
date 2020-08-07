package untrackthaturl

import (
	"fmt"
	"net/http"
	"net/url"

	"mvdan.cc/xurls/v2"
)

func fetchRedirectURL(resp *http.Response) string {
	location := resp.Header.Get("Location")
	if location != "" {
		return location
	}
	link := resp.Header.Get("Link")
	if link != "" {
		urls := xurls.Relaxed().FindAllString(link, -1)
		if len(urls) > 0 {
			return urls[0]
		}
	}

	return ""
}

func ResolveURL(u *url.URL) (*url.URL, error) {
	resp, err := http.Head(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return resp.Request.URL, nil
}
