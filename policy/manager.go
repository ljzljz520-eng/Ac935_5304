package policy

import (
	"errors"
	"fmt"
	"hospitaldesk/auth"
	"hospitaldesk/model"
	"hospitaldesk/storage"
	"strings"
	"time"
)

type Manager struct {
	store *storage.Store
	now   func() time.Time
}

func NewManager(store *storage.Store) *Manager {
	return &Manager{store: store, now: func() time.Time { return time.Date(2026, 1, 2, 9, 0, 0, 0, time.UTC) }}
}

func (m *Manager) CreateDraft(actor model.Employee, title, dept, content string) (model.PolicyDocument, error) {
	if err := auth.CanSubmitPolicy(actor); err != nil {
		return model.PolicyDocument{}, err
	}
	base := fmt.Sprintf("%s-%d", strings.ToLower(strings.ReplaceAll(title, " ", "-")), m.now().UnixNano())
	p := model.PolicyDocument{ID: base, Title: title, Department: dept, Content: content, Status: model.PolicyDraft, Version: 1, CreatedBy: actor.ID, CreatedAt: m.now(), UpdatedAt: m.now()}
	if err := model.ValidatePolicy(p); err != nil {
		return p, err
	}
	if err := m.store.SavePolicy(p); err != nil {
		return p, err
	}
	return p, nil
}

func (m *Manager) UpdateDraft(actor model.Employee, id, title, content string) (model.PolicyDocument, error) {
	p, err := m.store.GetPolicy(id)
	if err != nil {
		return p, err
	}
	if err := auth.CanSubmitPolicy(actor); err != nil {
		return p, err
	}
	if p.Status != model.PolicyDraft {
		return p, errors.New("only drafts can be edited")
	}
	if p.CreatedBy != actor.ID && actor.Role != model.RoleDirector {
		return p, errors.New("draft owner required")
	}
	if title != "" {
		p.Title = title
	}
	if content != "" {
		p.Content = content
	}
	p.UpdatedAt = m.now()
	if err := model.ValidatePolicy(p); err != nil {
		return p, err
	}
	return p, m.store.SavePolicy(p)
}

func (m *Manager) SubmitForReview(actor model.Employee, id string) (model.PolicyDocument, error) {
	p, err := m.store.GetPolicy(id)
	if err != nil {
		return p, err
	}
	if err := auth.CanSubmitPolicy(actor); err != nil {
		return p, err
	}
	if p.Status != model.PolicyDraft {
		return p, errors.New("policy is not a draft")
	}
	p.Status = model.PolicyPending
	p.UpdatedAt = m.now()
	if err := m.store.SavePolicy(p); err != nil {
		return p, err
	}
	event := model.ReviewEvent{ID: fmt.Sprintf("review-%s-submit", p.ID), PolicyID: p.ID, ActorID: actor.ID, Action: "submitted", At: m.now()}
	return p, m.store.SaveReview(event)
}

func (m *Manager) Approve(actor model.Employee, id string, effectiveFrom, effectiveUntil time.Time, comment string) (model.PolicyDocument, error) {
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
	p.Status = model.PolicyPublished
	p.ApprovedBy = actor.ID
	p.EffectiveFrom = effectiveFrom
	p.EffectiveUntil = effectiveUntil
	p.UpdatedAt = m.now()
	if err := model.ValidatePolicy(p); err != nil {
		return p, err
	}
	if err := m.store.SavePolicy(p); err != nil {
		return p, err
	}
	event := model.ReviewEvent{ID: fmt.Sprintf("review-%s-approve", p.ID), PolicyID: p.ID, ActorID: actor.ID, Action: "approved", At: m.now(), Comment: comment}
	return p, m.store.SaveReview(event)
}

func (m *Manager) Archive(actor model.Employee, id string) (model.PolicyDocument, error) {
	p, err := m.store.GetPolicy(id)
	if err != nil {
		return p, err
	}
	if err := auth.CanSubmitPolicy(actor); err != nil {
		return p, err
	}
	if p.Status != model.PolicyPublished {
		return p, errors.New("only published policies can be archived")
	}
	p.Status = model.PolicyArchived
	p.UpdatedAt = m.now()
	return p, m.store.SavePolicy(p)
}

func (m *Manager) GetForEmployee(actor model.Employee, id string) (model.PolicyDocument, error) {
	p, err := m.store.GetPolicy(id)
	if err != nil {
		return p, err
	}
	if err := auth.CanViewPolicy(actor, p, m.now()); err != nil {
		return p, err
	}
	return p, nil
}

func (m *Manager) Search(actor model.Employee, filter model.SearchFilter) ([]model.PolicyDocument, error) {
	if !actor.Active {
		return nil, auth.ErrInactiveEmployee
	}
	if filter.Department != "" && !actor.CanViewDepartment(filter.Department) {
		return nil, auth.ErrDepartmentDenied
	}
	items, err := m.store.SearchPolicies(filter)
	if err != nil {
		return nil, err
	}
	visible := make([]model.PolicyDocument, 0, len(items))
	for _, p := range items {
		if auth.CanViewPolicy(actor, p, m.now()) == nil {
			visible = append(visible, p)
		}
	}
	return visible, nil
}

func (m *Manager) ReviewHistory(id string) ([]model.ReviewEvent, error) {
	return m.store.PolicyReviews(id)
}

func (m *Manager) SetClock(now func() time.Time) error {
	if now == nil {
		return errors.New("clock cannot be nil")
	}
	m.now = now
	return nil
}
