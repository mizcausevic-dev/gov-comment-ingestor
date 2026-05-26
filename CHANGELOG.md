# Changelog

## v1.0.0-prod — 2026-05-26

- Platform/SRE hardening pass (Claude Code lane): added Dependabot (gomod + github-actions), Contributor Covenant code of conduct, static `/api/health.json` operability endpoint, and explicit Pages enablement on `configure-pages`.
- Verified production gates: `gofmt`, `go vet`, `go test ./...`, `cmd/prerender` + `cmd/demo` + `cmd/smoke` all green; AGPL-3.0-or-later, SECURITY.md, Pages-deployed SEO assets (robots/sitemap/OG).
- Deployed to GitHub Pages at `ingestor.kineticgain.com` (TLS).

## v0.1.0 - 2026-05-25

- initial public ship of the government comment ingestion control plane
- modeled adapter health, docket intake, normalization backlog, and evidence lineage
- prerendered static Pages deployment for `ingestor.kineticgain.com`
- Go validation lane with tests, demo, smoke, screenshots, and doctrine docs

