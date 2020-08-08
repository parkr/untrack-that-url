package untrackthaturl

import (
	"fmt"
	"net/http"
	"net/url"
)

// ResolveURL takes an input URL, sends a HEAD request to it and follows all HTTP redirects, returning the resulting URL.
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
