// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"os"
	"path/filepath"

	"github.com/mizcausevic-dev/gov-comment-ingestor/internal/service"
	"github.com/mizcausevic-dev/gov-comment-ingestor/internal/site"
)

func main() {
	must(os.RemoveAll("site"))
	must(os.MkdirAll(filepath.Join("site", "api", "dashboard"), 0o755))

	write("site/index.html", site.RenderOverview())
	write("site/ingest-lane.html", site.RenderIngestLane())
	write("site/source-adapters.html", site.RenderSourceAdapters())
	write("site/verification.html", site.RenderVerification())
	write("site/docs.html", site.RenderDocs())

	write("site/api/dashboard/summary/index.json", site.JSON(service.SummaryData()))
	write("site/api/ingest-lane.json", site.JSON(service.IngestLaneData()))
	write("site/api/source-adapters.json", site.JSON(service.SourceAdaptersData()))
	write("site/api/verification.json", site.JSON(service.VerificationData()))
	write("site/api/sample.json", site.JSON(service.SamplePayload()))

	if cname, err := os.ReadFile("CNAME"); err == nil {
		write("site/CNAME", string(cname))
	}
}

func write(path string, content string) {
	must(os.MkdirAll(filepath.Dir(path), 0o755))
	must(os.WriteFile(path, []byte(content), 0o644))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
