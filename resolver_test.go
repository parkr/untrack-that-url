package untrackthaturl

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestResolveURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "/final", http.StatusFound)
			return
		}
		if r.URL.Path == "/final" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.URL.Path == "/error" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}))
	defer ts.Close()

	t.Run("successful resolution with redirect", func(t *testing.T) {
		u, _ := url.Parse(ts.URL + "/redirect")
		res, err := ResolveURL(u)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expectedFinal := ts.URL + "/final"
		if res.URL != expectedFinal {
			t.Errorf("expected final URL %s, got %s", expectedFinal, res.URL)
		}

		if len(res.Trail) != 2 {
			t.Errorf("expected trail length 2, got %d", len(res.Trail))
		}
		if res.Trail[0] != ts.URL+"/redirect" {
			t.Errorf("expected trail[0] to be redirect URL, got %s", res.Trail[0])
		}
		if res.Trail[1] != expectedFinal {
			t.Errorf("expected trail[1] to be final URL, got %s", res.Trail[1])
		}
	})

	t.Run("resolution with error", func(t *testing.T) {
		u, _ := url.Parse(ts.URL + "/error")
		_, err := ResolveURL(u)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		expectedErr := fmt.Sprintf("unexpected status code 500 from %s/error", ts.URL)
		if err.Error() != expectedErr {
			t.Errorf("expected error %q, got %q", expectedErr, err.Error())
		}
	})
}
