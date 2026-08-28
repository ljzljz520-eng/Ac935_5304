package hospitaldesk

import (
	"hospitaldesk/audit"
	"hospitaldesk/model"
	"testing"
	"time"
)

func TestAuditSummary(t *testing.T) {
	records := []model.DownloadRecord{{ResourceType: "policy", Allowed: true, DownloadedAt: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)}, {ResourceType: "schedule", Allowed: false}}
	summary := audit.BuildSummary(records)
	if summary.Total != 2 || summary.Allowed != 1 || summary.Denied != 1 {
		t.Fatalf("summary=%+v", summary)
	}
	if len(audit.Recent(records, 1)) != 1 {
		t.Fatal("recent limit failed")
	}
}
