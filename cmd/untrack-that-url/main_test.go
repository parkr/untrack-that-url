package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRootHandler(t *testing.T) {
	server := httptest.NewServer(newServerMux())
	defer server.Close()

	t.Run("GET renders html", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/")
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("non-GET is rejected", func(t *testing.T) {
		resp, err := http.Post(server.URL+"/", "text/plain", strings.NewReader("nope"))
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
		}
	})
}

func TestResolveJSONHandler(t *testing.T) {
	redirectTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/redirect":
			http.Redirect(w, r, "/final", http.StatusFound)
		case "/final":
			w.WriteHeader(http.StatusOK)
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer redirectTarget.Close()

	server := httptest.NewServer(newServerMux())
	defer server.Close()

	t.Run("requires POST", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/resolve.json")
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
		}
		if got := resp.Header.Get("Content-Type"); got != "application/json; charset=utf-8" {
			t.Fatalf("unexpected content type: %q", got)
		}

		var body structuredError
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("unable to decode response: %v", err)
		}
		if body.Error != "unacceptable HTTP method" {
			t.Fatalf("unexpected error message: %q", body.Error)
		}
	})

	t.Run("requires URL parameter", func(t *testing.T) {
		resp, err := http.Post(server.URL+"/resolve.json", "application/x-www-form-urlencoded", strings.NewReader(""))
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}

		var body structuredError
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("unable to decode response: %v", err)
		}
		if body.Error != "url parameter required" {
			t.Fatalf("unexpected error message: %q", body.Error)
		}
	})

	t.Run("returns resolution info on success", func(t *testing.T) {
		values := url.Values{}
		values.Set("url", redirectTarget.URL+"/redirect")
		resp, err := http.Post(server.URL+"/resolve.json", "application/x-www-form-urlencoded", strings.NewReader(values.Encode()))
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
		if got := resp.Header.Get("Content-Type"); got != "application/json; charset=utf-8" {
			t.Fatalf("unexpected content type: %q", got)
		}

		var body urlResponse
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("unable to decode response: %v", err)
		}

		expectedURL := redirectTarget.URL + "/final"
		if body.URL != expectedURL {
			t.Fatalf("expected final url %q, got %q", expectedURL, body.URL)
		}
		if len(body.Trail) != 2 {
			t.Fatalf("expected 2 trail entries, got %d", len(body.Trail))
		}
		if body.Error != "" {
			t.Fatalf("expected empty error, got %q", body.Error)
		}
	})

	t.Run("returns trail and error on resolver failure", func(t *testing.T) {
		values := url.Values{}
		values.Set("url", redirectTarget.URL+"/error")
		resp, err := http.Post(server.URL+"/resolve.json", "application/x-www-form-urlencoded", strings.NewReader(values.Encode()))
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
		}

		var body urlResponse
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("unable to decode response: %v", err)
		}
		if !strings.Contains(body.Error, "unable to resolve url: unexpected status code 500") {
			t.Fatalf("unexpected error message: %q", body.Error)
		}
		if len(body.Trail) != 1 || body.Trail[0] != redirectTarget.URL+"/error" {
			t.Fatalf("unexpected trail: %#v", body.Trail)
		}
	})
}
