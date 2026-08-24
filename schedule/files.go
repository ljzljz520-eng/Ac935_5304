package schedule

import (
	"errors"
	"fmt"
	"hospitaldesk/auth"
	"hospitaldesk/model"
	"hospitaldesk/storage"
	"strings"
	"time"
)

type Service struct {
	store *storage.Store
	now   func() time.Time
}

func NewService(store *storage.Store) *Service {
	return &Service{store: store, now: func() time.Time { return time.Date(2026, 1, 2, 9, 0, 0, 0, time.UTC) }}
}

func (s *Service) Create(actor model.Employee, dept, period, content string) (model.ScheduleFile, error) {
	if err := auth.CanPublishSchedule(actor); err != nil {
		return model.ScheduleFile{}, err
	}
	f := model.ScheduleFile{ID: fmt.Sprintf("schedule-%s-%s", dept, strings.ReplaceAll(period, " ", "-")), Department: dept, Period: period, Content: content, Status: model.PolicyDraft, CreatedBy: actor.ID}
	if err := model.ValidateScheduleFile(f); err != nil {
		return f, err
	}
	return f, s.store.SaveSchedule(f)
}

func (s *Service) Publish(actor model.Employee, id string) (model.ScheduleFile, error) {
	f, err := s.store.GetSchedule(id)
	if err != nil {
		return f, err
	}
	if err := auth.CanPublishSchedule(actor); err != nil {
		return f, err
	}
	if f.Status != model.PolicyDraft {
		return f, errors.New("schedule is not a draft")
	}
	f.Status = model.PolicyPublished
	f.PublishedAt = s.now()
	return f, s.store.SaveSchedule(f)
}

func (s *Service) View(actor model.Employee, id string) (model.ScheduleFile, error) {
	f, err := s.store.GetSchedule(id)
	if err != nil {
		return f, err
	}
	if err := auth.CanViewSchedule(actor, f); err != nil {
		return f, err
	}
	return f, nil
}

func (s *Service) Search(actor model.Employee, dept, period string) ([]model.ScheduleFile, error) {
	if !actor.Active {
		return nil, auth.ErrInactiveEmployee
	}
	if dept != "" && !actor.CanViewDepartment(dept) {
		return nil, auth.ErrDepartmentDenied
	}
	all, err := s.store.ListSchedules()
	if err != nil {
		return nil, err
	}
	result := make([]model.ScheduleFile, 0)
	for _, f := range all {
		if f.Department == dept || dept == "" {
			if period == "" || f.Period == period {
				if auth.CanViewSchedule(actor, f) == nil {
					result = append(result, f)
				}
			}
		}
	}
	return result, nil
}

func (s *Service) SetClock(now func() time.Time) error {
	if now == nil {
		return errors.New("clock cannot be nil")
	}
	s.now = now
	return nil
}
