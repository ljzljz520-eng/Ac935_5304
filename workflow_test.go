package hospitaldesk

import (
	"errors"
	"hospitaldesk/model"
	"hospitaldesk/service"
	"hospitaldesk/storage"
	"testing"
	"time"
)

func testDesk(t *testing.T) *service.Desk {
	t.Helper()
	store, err := storage.Open(t.TempDir() + "/clinic.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	return service.New(store)
}

func staff() (model.Employee, model.Employee) {
	return model.Employee{ID: "director-1", Name: "主任", Department: "cardiology", Role: model.RoleDirector, Active: true}, model.Employee{ID: "employee-1", Name: "护士", Department: "cardiology", Role: model.RoleEmployee, Active: true}
}

func TestWorkflowPolicyApprovalAndDownload(t *testing.T) {
	desk := testDesk(t)
	director, employee := staff()
	if err := desk.RegisterEmployee(director); err != nil {
		t.Fatal(err)
	}
	if err := desk.RegisterEmployee(employee); err != nil {
		t.Fatal(err)
	}
	draft, err := desk.Policies.CreateDraft(director, "感染控制", "cardiology", "洗手和隔离要求")
	if err != nil {
		t.Fatal(err)
	}
	pending, err := desk.Policies.SubmitForReview(director, draft.ID)
	if err != nil {
		t.Fatal(err)
	}
	if pending.Status != model.PolicyPending {
		t.Fatalf("status=%s", pending.Status)
	}
	published, err := desk.Policies.Approve(director, draft.ID, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC), "ok")
	if err != nil {
		t.Fatal(err)
	}
	if published.Status != model.PolicyPublished {
		t.Fatalf("status=%s", published.Status)
	}
	if _, err := desk.RecordPolicyDownload(employee, published.ID); err != nil {
		t.Fatal(err)
	}
}

func TestWorkflowTrainingCompletion(t *testing.T) {
	desk := testDesk(t)
	director, employee := staff()
	if err := desk.RegisterEmployee(director); err != nil {
		t.Fatal(err)
	}
	if err := desk.RegisterEmployee(employee); err != nil {
		t.Fatal(err)
	}
	published, err := desk.PublishPolicyWorkflow(director, "用药规范", "cardiology", "双人核对", time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := desk.CompleteTraining(employee, employee.ID, published.ID, 96, "主任"); err != nil {
		t.Fatal(err)
	}
	records, err := desk.Training.ForEmployee(employee, employee.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].Score != 96 {
		t.Fatalf("records=%+v", records)
	}
	rate, err := desk.Training.CompletionRate(employee.ID, []string{published.ID})
	if err != nil {
		t.Fatal(err)
	}
	if rate != 1 {
		t.Fatalf("rate=%v", rate)
	}
}

func TestWorkflowScheduleAndDashboard(t *testing.T) {
	desk := testDesk(t)
	director, employee := staff()
	if err := desk.RegisterEmployee(director); err != nil {
		t.Fatal(err)
	}
	if err := desk.RegisterEmployee(employee); err != nil {
		t.Fatal(err)
	}
	file, err := desk.PublishScheduleWorkflow(director, employee.Department, "2026-W02", "周一早班")
	if err != nil {
		t.Fatal(err)
	}
	if file.Status != model.PolicyPublished {
		t.Fatalf("status=%s", file.Status)
	}
	viewed, err := desk.ViewSchedule(employee, file.ID)
	if err != nil {
		t.Fatal(err)
	}
	if viewed.Content != "周一早班" {
		t.Fatalf("content=%s", viewed.Content)
	}
	dashboard, err := desk.Dashboard(employee)
	if err != nil {
		t.Fatal(err)
	}
	if len(dashboard.Schedules) != 1 {
		t.Fatalf("schedules=%d", len(dashboard.Schedules))
	}
	if errors.Is(err, nil) == false {
		t.Fatal("unexpected error")
	}
}

func TestEmployeeCannotDownloadUnpublishedPolicy(t *testing.T) {
	desk := testDesk(t)
	director, employee := staff()
	if err := desk.RegisterEmployee(director); err != nil {
		t.Fatal(err)
	}
	if err := desk.RegisterEmployee(employee); err != nil {
		t.Fatal(err)
	}
	draft, err := desk.Policies.CreateDraft(director, "未发布制度", employee.Department, "草案内容")
	if err != nil {
		t.Fatal(err)
	}
	link, err := desk.Shares.Create(director, draft)
	if err != nil {
		t.Fatal(err)
	}
	_, err = desk.RecordShareDownload(employee, link.Token)
	if err == nil {
		t.Fatalf("employee downloaded unpublished policy")
	}
}
