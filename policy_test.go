package hospitaldesk

import (
	"hospitaldesk/auth"
	"hospitaldesk/model"
	"testing"
	"time"
)

func TestPolicyValidityWindow(t *testing.T) {
	p := model.PolicyDocument{Status: model.PolicyPublished, EffectiveFrom: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), EffectiveUntil: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)}
	if p.IsCurrentlyValid(time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)) == false {
		t.Fatal("expected valid")
	}
	if p.IsCurrentlyValid(time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("expected expired")
	}
}

func TestShareExpiryIsDenied(t *testing.T) {
	actor := model.Employee{ID: "e", Department: "x", Role: model.RoleEmployee, Active: true}
	p := model.PolicyDocument{ID: "p", Department: "x", Status: model.PolicyPublished}
	link := model.ShareLink{Token: "t", ResourceID: "p", ResourceType: "policy", Department: "x", ExpiresAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
	if err := auth.ValidateShare(actor, link, p, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)); err != auth.ErrShareInvalid {
		t.Fatalf("err=%v", err)
	}
}
