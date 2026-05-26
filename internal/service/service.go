// SPDX-License-Identifier: AGPL-3.0-or-later

package service

type Summary struct {
	AdaptersHealthy   int     `json:"adaptersHealthy"`
	AdaptersDegraded  int     `json:"adaptersDegraded"`
	Comments24h       int     `json:"comments24h"`
	NormalizationRate float64 `json:"normalizationRate"`
	ThreadingFlags    int     `json:"threadingFlags"`
	EvidenceGaps      int     `json:"evidenceGaps"`
}

type IngestItem struct {
	Queue          string `json:"queue"`
	Agency         string `json:"agency"`
	Window         string `json:"window"`
	Status         string `json:"status"`
	NewComments    int    `json:"newComments"`
	DriftSignal    string `json:"driftSignal"`
	AttachmentRisk string `json:"attachmentRisk"`
	Owner          string `json:"owner"`
	NextAction     string `json:"nextAction"`
}

type SourceAdapter struct {
	Name             string `json:"name"`
	Feed             string `json:"feed"`
	Status           string `json:"status"`
	Lag              string `json:"lag"`
	ParserVersion    string `json:"parserVersion"`
	CommentCoverage  string `json:"commentCoverage"`
	AttachmentPolicy string `json:"attachmentPolicy"`
	Escalation       string `json:"escalation"`
}

type VerificationGate struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Evidence   string `json:"evidence"`
	Risk       string `json:"risk"`
	NextAction string `json:"nextAction"`
}

type Payload struct {
	Title           string             `json:"title"`
	Summary         Summary            `json:"summary"`
	IngestLane      []IngestItem       `json:"ingestLane"`
	SourceAdapters  []SourceAdapter    `json:"sourceAdapters"`
	Verification    []VerificationGate `json:"verification"`
	SignalNarrative []string           `json:"signalNarrative"`
}

func SamplePayload() Payload {
	return Payload{
		Title: "Gov Comment Ingestor",
		Summary: Summary{
			AdaptersHealthy:   4,
			AdaptersDegraded:  2,
			Comments24h:       1284,
			NormalizationRate: 96.8,
			ThreadingFlags:    14,
			EvidenceGaps:      3,
		},
		IngestLane: []IngestItem{
			{
				Queue:          "EPA methane reporting amendments",
				Agency:         "EPA",
				Window:         "18h to analyst handoff",
				Status:         "Watch",
				NewComments:    186,
				DriftSignal:    "Attachment filenames changed after midnight refresh",
				AttachmentRisk: "2 PDF exhibits still unparsed",
				Owner:          "Policy ops",
				NextAction:     "Re-run attachment parser before legal packet lock",
			},
			{
				Queue:          "DOT autonomous freight waiver docket",
				Agency:         "DOT",
				Window:         "9h to executive review",
				Status:         "Critical",
				NewComments:    92,
				DriftSignal:    "Thread IDs doubled after source replay",
				AttachmentRisk: "Narrative references mis-threaded",
				Owner:          "Gov affairs",
				NextAction:     "Deduplicate thread graph before briefing bundle ships",
			},
			{
				Queue:          "FTC AI disclosure rulemaking",
				Agency:         "FTC",
				Window:         "36h to synthesis package",
				Status:         "Healthy",
				NewComments:    240,
				DriftSignal:    "Stable",
				AttachmentRisk: "No missing exhibits",
				Owner:          "Research ops",
				NextAction:     "Continue nightly ingest cadence",
			},
			{
				Queue:          "State utility resilience comments",
				Agency:         "MA DPU",
				Window:         "22h to counsel review",
				Status:         "Watch",
				NewComments:    71,
				DriftSignal:    "Source HTML wrapped quote blocks differently",
				AttachmentRisk: "1 spreadsheet attachment pending checksum",
				Owner:          "Infrastructure counsel",
				NextAction:     "Normalize quotes and confirm spreadsheet lineage",
			},
		},
		SourceAdapters: []SourceAdapter{
			{
				Name:             "regulations-dot-gov",
				Feed:             "API pull + attachment manifest",
				Status:           "Healthy",
				Lag:              "6m",
				ParserVersion:    "go-regtext-1.4.2",
				CommentCoverage:  "100% root threads",
				AttachmentPolicy: "hash + mime + OCR",
				Escalation:       "None",
			},
			{
				Name:             "federal-register-ingest",
				Feed:             "JSON notices + linked PDFs",
				Status:           "Healthy",
				Lag:              "11m",
				ParserVersion:    "fr-parser-2.0.1",
				CommentCoverage:  "98.9%",
				AttachmentPolicy: "PDF parse + cite map",
				Escalation:       "Monitor OCR backlog",
			},
			{
				Name:             "epa-exhibits-dropbox",
				Feed:             "Manual export landing zone",
				Status:           "Watch",
				Lag:              "47m",
				ParserVersion:    "attachment-normalizer-0.9.7",
				CommentCoverage:  "95.1%",
				AttachmentPolicy: "checksum + filename lineage",
				Escalation:       "Exhibit naming drift after contractor handoff",
			},
			{
				Name:             "state-commission-html-scraper",
				Feed:             "HTML scrape + thread reconstruction",
				Status:           "Critical",
				Lag:              "2h 08m",
				ParserVersion:    "commission-scrape-0.6.4",
				CommentCoverage:  "89.4%",
				AttachmentPolicy: "HTML + CSV attachment mirror",
				Escalation:       "DOM shift broke quote + attachment binding",
			},
			{
				Name:             "sec-comment-mirror",
				Feed:             "RSS + archive PDF fetch",
				Status:           "Healthy",
				Lag:              "14m",
				ParserVersion:    "sec-ingest-1.1.0",
				CommentCoverage:  "97.6%",
				AttachmentPolicy: "archive mirror + sha256",
				Escalation:       "None",
			},
			{
				Name:             "counsel-inbox-forwarder",
				Feed:             "Email attachment intake",
				Status:           "Watch",
				Lag:              "31m",
				ParserVersion:    "mail-bridge-0.8.3",
				CommentCoverage:  "94.8%",
				AttachmentPolicy: "mail header + attachment lineage",
				Escalation:       "Missing sender alias map for two agencies",
			},
		},
		Verification: []VerificationGate{
			{
				Name:       "adapter health matrix",
				Status:     "watch",
				Evidence:   "2 of 6 adapters degraded; DOT and state utility queues still ingesting",
				Risk:       "review packets can ship stale thread references",
				NextAction: "stabilize scraper + replay queue before next analyst pack",
			},
			{
				Name:       "evidence lineage",
				Status:     "healthy",
				Evidence:   "sha256 lineage and attachment mirrors present for 96.8% of comments",
				Risk:       "low",
				NextAction: "close remaining exhibit OCR backlog",
			},
			{
				Name:       "duplicate thread suppression",
				Status:     "critical",
				Evidence:   "14 thread collisions after replay on DOT waiver docket",
				Risk:       "briefing narratives can double-count opposition themes",
				NextAction: "rebuild thread map before executive review",
			},
			{
				Name:       "submission-safe export",
				Status:     "healthy",
				Evidence:   "CSV + JSON export schemas match downstream synthesis contract",
				Risk:       "low",
				NextAction: "continue nightly contract check",
			},
		},
		SignalNarrative: []string{
			"Stable ingestion matters more than raw volume when legal and policy teams depend on evidence traceability.",
			"Parser drift is a release blocker when quote threading and attachment lineage can change interpretation.",
			"Buyer-safe operator tooling makes source health visible before review and submission workflows inherit bad data.",
		},
	}
}

func SummaryData() Summary {
	return SamplePayload().Summary
}

func IngestLaneData() []IngestItem {
	return SamplePayload().IngestLane
}

func SourceAdaptersData() []SourceAdapter {
	return SamplePayload().SourceAdapters
}

func VerificationData() []VerificationGate {
	return SamplePayload().Verification
}
