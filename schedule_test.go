package hospitaldesk

import (
	"hospitaldesk/model"
	"hospitaldesk/schedule"
	"testing"
)

func TestScheduleGrouping(t *testing.T) {
	files := []model.ScheduleFile{{ID: "2", Period: "w1", Department: "z"}, {ID: "1", Period: "w1", Department: "a"}, {ID: "3", Period: "w2", Department: "a"}}
	groups := schedule.GroupByPeriod(files)
	if len(groups["w1"]) != 2 || groups["w1"][0].Department != "a" {
		t.Fatalf("groups=%v", groups)
	}
	if !schedule.IsCurrentPeriod(files[0], " W1 ") {
		t.Fatal("period mismatch")
	}
}
