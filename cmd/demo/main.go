package main

import (
	"fmt"

	"github.com/mizcausevic-dev/gov-comment-ingestor/internal/service"
)

func main() {
	summary := service.SummaryData()
	fmt.Printf("healthy_adapters=%d degraded_adapters=%d comments_24h=%d normalization_rate=%.1f threading_flags=%d evidence_gaps=%d\n",
		summary.AdaptersHealthy,
		summary.AdaptersDegraded,
		summary.Comments24h,
		summary.NormalizationRate,
		summary.ThreadingFlags,
		summary.EvidenceGaps,
	)
}
