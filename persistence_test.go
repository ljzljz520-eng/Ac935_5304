package hospitaldesk

import (
	"hospitaldesk/model"
	"hospitaldesk/storage"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/persist.db"
	store, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	items := []any{model.Employee{ID: "e1", Name: "A", Department: "x", Role: model.RoleEmployee, Active: true}, model.PolicyDocument{ID: "p1", Title: "P", Department: "x", Content: "C", Status: model.PolicyPublished, Version: 1}, model.TrainingRecord{ID: "t1", EmployeeID: "e1", PolicyID: "p1", Score: 80}, model.ScheduleFile{ID: "s1", Department: "x", Period: "W1", Content: "C", Status: model.PolicyPublished}}
	if err := store.SaveEmployee(items[0].(model.Employee)); err != nil {
		t.Fatal(err)
	}
	if err := store.SavePolicy(items[1].(model.PolicyDocument)); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveTraining(items[2].(model.TrainingRecord)); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveSchedule(items[3].(model.ScheduleFile)); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	if _, err := reopened.GetEmployee("e1"); err != nil {
		t.Fatal(err)
	}
	if _, err := reopened.GetPolicy("p1"); err != nil {
		t.Fatal(err)
	}
	if _, err := reopened.GetTraining("t1"); err != nil {
		t.Fatal(err)
	}
	if _, err := reopened.GetSchedule("s1"); err != nil {
		t.Fatal(err)
	}
}
