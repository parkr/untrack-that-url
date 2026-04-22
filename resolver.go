package untrackthaturl

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	publicResolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * 5,
			}
			return d.DialContext(ctx, network, "1.1.1.1:53")
		},
	}

	httpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
				Resolver:  publicResolver,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 20 * time.Second,
	}
)

type Resolution struct {
	URL   string   `json:"url"`
	Trail []string `json:"trail"`
}

// ResolveURL takes an input URL, sends a HEAD request to it and follows all HTTP redirects, returning the resulting URL and the trail of URLs followed.
func ResolveURL(u *url.URL) (*Resolution, error) {
	res := &Resolution{
		Trail: []string{u.String()},
	}

	// Create a local copy to ensure thread-safety when setting CheckRedirect.
	client := *httpClient
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return fmt.Errorf("too many redirects")
		}
		res.Trail = append(res.Trail, req.URL.String())
		return nil
	}

	resp, err := client.Head(u.String())
	if err != nil {
		return res, fmt.Errorf("error during request to %s: %v", res.Trail[len(res.Trail)-1], err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return res, fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, resp.Request.URL.String())
	}

	res.URL = resp.Request.URL.String()
	return res, nil
}
