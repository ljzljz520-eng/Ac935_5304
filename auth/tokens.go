package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hospitaldesk/model"
	"time"
)

func CreateShare(actor model.Employee, resourceID, resourceType, department string, now time.Time) (model.ShareLink, error) {
	if !actor.Active {
		return model.ShareLink{}, ErrInactiveEmployee
	}
	if !actor.CanViewDepartment(department) && actor.Role != model.RoleDirector {
		return model.ShareLink{}, ErrDepartmentDenied
	}
	if resourceID == "" || resourceType == "" {
		return model.ShareLink{}, errors.New("share resource is required")
	}
	seed := fmt.Sprintf("%s:%s:%s", actor.ID, resourceID, now.UTC().Format(time.RFC3339))
	sum := sha256.Sum256([]byte(seed))
	return model.ShareLink{Token: hex.EncodeToString(sum[:12]), ResourceID: resourceID, ResourceType: resourceType, Department: department, CreatedBy: actor.ID, ExpiresAt: now.Add(24 * time.Hour)}, nil
}

func RevokeShare(actor model.Employee, link model.ShareLink) error {
	if err := CanSubmitPolicy(actor); err != nil {
		return err
	}
	if link.CreatedBy != actor.ID && actor.Role != model.RoleDirector {
		return errors.New("share owner required")
	}
	link.Revoked = true
	return nil
}

func ExplainAccess(err error) string {
	if err == nil {
		return "allowed"
	}
	if errors.Is(err, ErrUnpublished) {
		return "awaiting director approval"
	}
	if errors.Is(err, ErrExpired) {
		return "outside effective period"
	}
	if errors.Is(err, ErrDepartmentDenied) {
		return "department restriction"
	}
	return err.Error()
}
