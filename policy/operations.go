package policy

import (
	"errors"
	"fmt"
	"hospitaldesk/auth"
	"hospitaldesk/model"
	"sort"
	"strings"
	"time"
)

type BatchResult struct {
	Processed int
	Succeeded int
	Failed    int
	Errors    []string
}

func (m *Manager) CloneDraft(actor model.Employee, sourceID, title string) (model.PolicyDocument, error) {
	if err := auth.CanSubmitPolicy(actor); err != nil {
		return model.PolicyDocument{}, err
	}
	source, err := m.store.GetPolicy(sourceID)
	if err != nil {
		return model.PolicyDocument{}, err
	}
	if source.Status == model.PolicyArchived {
		return model.PolicyDocument{}, errors.New("archived policy cannot be cloned")
	}
	clone := source
	clone.ID = fmt.Sprintf("%s-v%d", source.ID, source.Version+1)
	clone.Version = source.Version + 1
	clone.Status = model.PolicyDraft
	clone.ApprovedBy = ""
	clone.CreatedBy = actor.ID
	clone.Title = title
	clone.CreatedAt = m.now()
	clone.UpdatedAt = m.now()
	if err := model.ValidatePolicy(clone); err != nil {
		return clone, err
	}
	return clone, m.store.SavePolicy(clone)
}

func (m *Manager) Reject(actor model.Employee, id, comment string) (model.PolicyDocument, error) {
	p, err := m.store.GetPolicy(id)
	if err != nil {
		return p, err
	}
	if err := auth.CanApprovePolicy(actor); err != nil {
		return p, err
	}
	if p.Status != model.PolicyPending {
		return p, errors.New("policy is not awaiting review")
	}
	p.Status = model.PolicyDraft
	p.UpdatedAt = m.now()
	if err := m.store.SavePolicy(p); err != nil {
		return p, err
	}
	event := model.ReviewEvent{ID: fmt.Sprintf("review-%s-reject-%d", p.ID, m.now().UnixNano()), PolicyID: p.ID, ActorID: actor.ID, Action: "rejected", At: m.now(), Comment: comment}
	return p, m.store.SaveReview(event)
}

func (m *Manager) BatchArchive(actor model.Employee, ids []string) BatchResult {
	result := BatchResult{Errors: make([]string, 0)}
	for _, id := range ids {
		result.Processed++
		if _, err := m.Archive(actor, id); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, id+": "+err.Error())
		} else {
			result.Succeeded++
		}
	}
	return result
}

func SortByVersion(policies []model.PolicyDocument) []model.PolicyDocument {
	result := append([]model.PolicyDocument(nil), policies...)
	sort.Slice(result, func(i, j int) bool {
		if result[i].Version == result[j].Version {
			return result[i].UpdatedAt.After(result[j].UpdatedAt)
		}
		return result[i].Version > result[j].Version
	})
	return result
}

func GroupVersions(policies []model.PolicyDocument) map[string][]model.PolicyDocument {
	groups := make(map[string][]model.PolicyDocument)
	for _, p := range policies {
		key := strings.TrimSuffix(p.ID, fmt.Sprintf("-v%d", p.Version))
		groups[key] = append(groups[key], p)
	}
	for key := range groups {
		groups[key] = SortByVersion(groups[key])
	}
	return groups
}

func ValidateApprovalWindow(from, until, now time.Time) error {
	if from.IsZero() {
		return errors.New("effective start is required")
	}
	if !until.IsZero() && until.Before(from) {
		return errors.New("effective end precedes start")
	}
	if until.IsZero() {
		return nil
	}
	if until.Before(now) {
		return errors.New("effective window is already expired")
	}
	return nil
}

func IsReviewable(policy model.PolicyDocument) bool {
	return policy.Status == model.PolicyPending && strings.TrimSpace(policy.Content) != ""
}

func AuditActionAt(actor model.Employee, action string, at time.Time) model.ReviewEvent {
	return model.ReviewEvent{ID: fmt.Sprintf("action-%s-%d", action, at.UnixNano()), ActorID: actor.ID, Action: action, At: at}
}
