# Architecture

## Control plane

`gov-comment-ingestor` models the intake tier in front of regulated workflow analysis:

1. source adapters pull docket payloads from federal, state, and sector-specific feeds
2. normalization maps those payloads into a review-safe event structure
3. evidence lineage ties each normalized comment back to source, thread, and attachment status
4. operators see lag, parser drift, duplicate threading, and missing evidence before downstream reviewers act on stale data

## Application shape

- `cmd/server`
  - local web server with HTML routes and JSON API endpoints
- `cmd/prerender`
  - writes static HTML and JSON artifacts into `site/` for GitHub Pages deploys
- `cmd/demo`
  - prints summary signals for release verification
- `cmd/smoke`
  - validates the prerendered site structure and key route content
- `internal/service`
  - pure modeled data and summary functions
- `internal/site`
  - HTML rendering helpers shared by server + prerender paths

## SEO/AEO posture

- crawlable HTML pages for all primary routes
- canonical URLs added at deploy time
- `robots.txt` + `sitemap.xml` emitted during Pages deploy
- operator-oriented copy that is still buyer-readable and answer-engine legible

