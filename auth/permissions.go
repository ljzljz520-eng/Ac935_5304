package auth

import (
	"errors"
	"hospitaldesk/model"
	"strings"
	"time"
)

var (
	ErrInactiveEmployee = errors.New("employee is inactive")
	ErrDepartmentDenied = errors.New("employee is outside department")
	ErrUnpublished      = errors.New("policy is not published")
	ErrExpired          = errors.New("policy is outside effective period")
	ErrShareInvalid     = errors.New("share link is invalid")
)

func CanSubmitPolicy(actor model.Employee) error {
	if !actor.Active {
		return ErrInactiveEmployee
	}
	if !actor.CanManagePolicies() {
		return errors.New("policy management role required")
	}
	return nil
}

func CanApprovePolicy(actor model.Employee) error {
	if err := CanSubmitPolicy(actor); err != nil {
		return err
	}
	if actor.Role != model.RoleDirector {
		return errors.New("director approval required")
	}
	return nil
}

func CanViewPolicy(actor model.Employee, p model.PolicyDocument, now time.Time) error {
	if !actor.Active {
		return ErrInactiveEmployee
	}
	if !actor.CanViewDepartment(p.Department) {
		return ErrDepartmentDenied
	}
	if p.Status != model.PolicyPublished {
		return ErrUnpublished
	}
	if (!p.EffectiveFrom.IsZero() && now.Before(p.EffectiveFrom)) || (!p.EffectiveUntil.IsZero() && now.After(p.EffectiveUntil)) {
		return ErrExpired
	}
	return nil
}

func CanViewSchedule(actor model.Employee, f model.ScheduleFile) error {
	if !actor.Active {
		return ErrInactiveEmployee
	}
	if !actor.CanViewDepartment(f.Department) {
		return ErrDepartmentDenied
	}
	if !f.IsPublished() {
		return ErrUnpublished
	}
	return nil
}

func ValidateShare(actor model.Employee, link model.ShareLink, p model.PolicyDocument, now time.Time) error {
	if link.Token == "" || link.Revoked || link.ResourceID != p.ID || link.ResourceType != "policy" {
		return ErrShareInvalid
	}
	if !link.ExpiresAt.IsZero() && now.After(link.ExpiresAt) {
		return ErrShareInvalid
	}
	if !actor.Active {
		return ErrInactiveEmployee
	}
	if actor.Role == model.RoleDirector || strings.EqualFold(link.Department, actor.Department) {
		return nil
	}
	if actor.Role == model.RoleEmployee && link.Department != "" {
		return nil
	}
	return ErrDepartmentDenied
}

func CanRecordTraining(actor model.Employee, employeeID string) error {
	if !actor.Active {
		return ErrInactiveEmployee
	}
	if actor.ID == employeeID || actor.CanManagePolicies() {
		return nil
	}
	return errors.New("training record access denied")
}

func CanPublishSchedule(actor model.Employee) error {
	if err := CanSubmitPolicy(actor); err != nil {
		return err
	}
	return nil
}
