package hospitaldesk

import (
	"hospitaldesk/model"
	"hospitaldesk/training"
	"testing"
)

func TestTrainingFilteringAndSummary(t *testing.T) {
	records := []model.TrainingRecord{{EmployeeID: "e1", PolicyID: "p1", Score: 80, Notes: "safe"}, {EmployeeID: "e2", PolicyID: "p2", Score: 40, Notes: "retry"}}
	filtered := training.FilterRecords(records, "safe", 60)
	if len(filtered) != 1 {
		t.Fatalf("len=%d", len(filtered))
	}
	passed, total := training.Summarize(records)
	if passed != 1 || total != 2 {
		t.Fatalf("%d/%d", passed, total)
	}
}
