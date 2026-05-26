package service

import "testing"

func TestSamplePayloadShape(t *testing.T) {
	payload := SamplePayload()

	if payload.Summary.AdaptersHealthy+payload.Summary.AdaptersDegraded != len(payload.SourceAdapters) {
		t.Fatalf("adapter count mismatch: summary=%d adapters=%d", payload.Summary.AdaptersHealthy+payload.Summary.AdaptersDegraded, len(payload.SourceAdapters))
	}

	if len(payload.IngestLane) < 4 {
		t.Fatalf("expected at least four ingest items, got %d", len(payload.IngestLane))
	}

	if payload.Verification[2].Status != "critical" {
		t.Fatalf("expected duplicate thread suppression to be critical")
	}
}
