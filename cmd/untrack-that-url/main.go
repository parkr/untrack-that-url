package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	untrackthaturl "github.com/parkr/untrack-that-url"
)

type urlResponse struct {
	URL   string   `json:"url,omitempty"`
	Trail []string `json:"trail,omitempty"`
	Error string   `json:"error,omitempty"`
}

type structuredError struct {
	Error string `json:"error"`
}

func jsonError(w http.ResponseWriter, errMessage string, code int) error {
	encoded, err := json.Marshal(structuredError{errMessage})
	if err != nil {
		return err
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	_, err = w.Write(encoded)
	return err
}

func main() {
	bind := flag.String("http", ":8080", "IP:PORT to bind http server to")
	flag.Parse()
	if *bind == "" {
		fmt.Println("fatal: -bind flag requires a value")
		os.Exit(1)
	}

	http.HandleFunc("/resolve.json", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			err := jsonError(w, "unacceptable HTTP method", http.StatusMethodNotAllowed)
			if err != nil {
				http.Error(w, "unexpected error: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}
		stringURL := req.FormValue("url")
		if stringURL == "" {
			err := jsonError(w, "url parameter required", http.StatusBadRequest)
			if err != nil {
				http.Error(w, "unexpected error: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}
		u, err := url.Parse(stringURL)
		if err != nil {
			err := jsonError(w, "url parameter invalid: "+err.Error(), http.StatusBadRequest)
			if err != nil {
				http.Error(w, "unexpected error: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf8")

		resolution, err := untrackthaturl.ResolveURL(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(urlResponse{
				Trail: resolution.Trail,
				Error: "unable to resolve url: " + err.Error(),
			})
			return
		}

		err = json.NewEncoder(w).Encode(urlResponse{
			URL:   resolution.URL,
			Trail: resolution.Trail,
		})
		if err != nil {
			err := jsonError(w, "unable to encode json: "+err.Error(), http.StatusInternalServerError)
			if err != nil {
				http.Error(w, "unexpected error: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, "unacceptable method", http.StatusMethodNotAllowed)
			return
		}
		if err := untrackthaturl.RenderHTML(w); err != nil {
			http.Error(w, "error rendering html", http.StatusInternalServerError)
		}
	})

	fmt.Printf("untrack-that-url: serving on %q\n", *bind)
	if err := http.ListenAndServe(*bind, nil); err != nil {
		fmt.Printf("fatal: error serving on %q: %v", *bind, err)
		os.Exit(1)
	}
}
