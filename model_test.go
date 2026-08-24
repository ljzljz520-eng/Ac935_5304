package hospitaldesk

import (
	"hospitaldesk/model"
	"testing"
)

func TestModelValidation(t *testing.T) {
	if err := model.ValidateEmployee(model.Employee{}); err == nil {
		t.Fatal("expected employee validation")
	}
	if err := model.ValidatePolicy(model.PolicyDocument{}); err == nil {
		t.Fatal("expected policy validation")
	}
	if err := model.ValidateTrainingRecord(model.TrainingRecord{ID: "t", EmployeeID: "e", PolicyID: "p", Score: 101}); err == nil {
		t.Fatal("expected score validation")
	}
}
