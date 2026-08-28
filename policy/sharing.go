package policy

import (
	"errors"
	"hospitaldesk/auth"
	"hospitaldesk/model"
	"hospitaldesk/storage"
	"time"
)

type ShareService struct {
	store *storage.Store
	now   func() time.Time
}

func NewShareService(store *storage.Store) *ShareService {
	return &ShareService{store: store, now: func() time.Time { return time.Date(2026, 1, 2, 9, 0, 0, 0, time.UTC) }}
}

func (s *ShareService) Create(actor model.Employee, policy model.PolicyDocument) (model.ShareLink, error) {
	link, err := auth.CreateShare(actor, policy.ID, "policy", policy.Department, s.now())
	if err != nil {
		return link, err
	}
	return link, s.store.SaveShare(link)
}

func (s *ShareService) Download(actor model.Employee, token string) (model.PolicyDocument, error) {
	link, err := s.store.GetShare(token)
	if err != nil {
		return model.PolicyDocument{}, err
	}
	p, err := s.store.GetPolicy(link.ResourceID)
	if err != nil {
		return p, err
	}
	if link.Revoked {
		return p, auth.ErrShareInvalid
	}
	if err := auth.ValidateShare(actor, link, p, s.now()); err != nil {
		return p, err
	}
	return p, nil
}

func (s *ShareService) Revoke(actor model.Employee, token string) error {
	link, err := s.store.GetShare(token)
	if err != nil {
		return err
	}
	if err := auth.RevokeShare(actor, link); err != nil {
		return err
	}
	link.Revoked = true
	return s.store.SaveShare(link)
}

func (s *ShareService) SetClock(now func() time.Time) error {
	if now == nil {
		return errors.New("clock cannot be nil")
	}
	s.now = now
	return nil
}
