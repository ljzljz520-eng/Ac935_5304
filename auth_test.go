package hospitaldesk

import (
	"hospitaldesk/auth"
	"hospitaldesk/model"
	"testing"
	"time"
)

func TestDepartmentAuthorization(t *testing.T) {
	actor := model.Employee{ID: "e", Department: "x", Role: model.RoleEmployee, Active: true}
	p := model.PolicyDocument{ID: "p", Department: "y", Status: model.PolicyPublished}
	if err := auth.CanViewPolicy(actor, p, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)); err != auth.ErrDepartmentDenied {
		t.Fatalf("err=%v", err)
	}
}
