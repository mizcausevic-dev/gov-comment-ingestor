// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/mizcausevic-dev/gov-comment-ingestor/internal/service"
	"github.com/mizcausevic-dev/gov-comment-ingestor/internal/site"
)

func main() {
	port := 5514
	if raw := os.Getenv("PORT"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil {
			port = parsed
		}
	}
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			html(w, site.RenderOverview())
		case "/ingest-lane":
			html(w, site.RenderIngestLane())
		case "/source-adapters":
			html(w, site.RenderSourceAdapters())
		case "/verification":
			html(w, site.RenderVerification())
		case "/docs":
			html(w, site.RenderDocs())
		case "/api/dashboard/summary":
			json(w, service.SummaryData())
		case "/api/ingest-lane":
			json(w, service.IngestLaneData())
		case "/api/source-adapters":
			json(w, service.SourceAdaptersData())
		case "/api/verification":
			json(w, service.VerificationData())
		case "/api/sample":
			json(w, service.SamplePayload())
		default:
			http.NotFound(w, r)
		}
	})

	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Gov Comment Ingestor listening on http://%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}

func html(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(body))
}

func json(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write([]byte(site.JSON(payload)))
}
