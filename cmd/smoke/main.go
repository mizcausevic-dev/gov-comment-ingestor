package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	required := map[string]string{
		filepath.Join("site", "index.html"):              "Government comment feeds fail",
		filepath.Join("site", "ingest-lane.html"):        "Ingest lane",
		filepath.Join("site", "source-adapters.html"):    "Source adapters",
		filepath.Join("site", "verification.html"):       "Verification",
		filepath.Join("site", "docs.html"):               "The ingestion tier is part of the product surface",
		filepath.Join("site", "api", "sample.json"):      "Gov Comment Ingestor",
		filepath.Join("site", "api", "ingest-lane.json"): "DOT autonomous freight waiver docket",
	}

	for path, needle := range required {
		data, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		if !bytes.Contains(data, []byte(needle)) {
			panic(fmt.Sprintf("%s missing %q", path, needle))
		}
	}

	fmt.Println("smoke ok")
}
