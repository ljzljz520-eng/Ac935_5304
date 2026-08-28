package model

import (
	"errors"
	"strings"
)

func ValidateEmployee(e Employee) error {
	if strings.TrimSpace(e.ID) == "" {
		return errors.New("employee id is required")
	}
	if strings.TrimSpace(e.Name) == "" {
		return errors.New("employee name is required")
	}
	if strings.TrimSpace(e.Department) == "" {
		return errors.New("department is required")
	}
	if e.Role != RoleDirector && e.Role != RoleSupervisor && e.Role != RoleEmployee {
		return errors.New("unknown employee role")
	}
	return nil
}

func ValidatePolicy(p PolicyDocument) error {
	if strings.TrimSpace(p.ID) == "" {
		return errors.New("policy id is required")
	}
	if strings.TrimSpace(p.Title) == "" {
		return errors.New("policy title is required")
	}
	if strings.TrimSpace(p.Department) == "" {
		return errors.New("policy department is required")
	}
	if strings.TrimSpace(p.Content) == "" {
		return errors.New("policy content is required")
	}
	if p.Version < 1 {
		return errors.New("policy version must be positive")
	}
	if !p.EffectiveUntil.IsZero() && !p.EffectiveFrom.IsZero() && p.EffectiveUntil.Before(p.EffectiveFrom) {
		return errors.New("invalid policy validity window")
	}
	return nil
}

func ValidateTrainingRecord(r TrainingRecord) error {
	if r.ID == "" || r.EmployeeID == "" || r.PolicyID == "" {
		return errors.New("training identifiers are required")
	}
	if r.Score < 0 || r.Score > 100 {
		return errors.New("training score must be between 0 and 100")
	}
	return nil
}

func ValidateScheduleFile(s ScheduleFile) error {
	if s.ID == "" || s.Department == "" || s.Period == "" {
		return errors.New("schedule identity is required")
	}
	if s.Content == "" {
		return errors.New("schedule content is required")
	}
	return nil
}
