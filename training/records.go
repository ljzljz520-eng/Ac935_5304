package training

import (
	"errors"
	"fmt"
	"hospitaldesk/auth"
	"hospitaldesk/model"
	"hospitaldesk/storage"
	"time"
)

type Service struct {
	store *storage.Store
	now   func() time.Time
}

func NewService(store *storage.Store) *Service {
	return &Service{store: store, now: func() time.Time { return time.Date(2026, 1, 2, 9, 0, 0, 0, time.UTC) }}
}

func (s *Service) Record(actor model.Employee, employeeID, policyID string, score int, trainer, notes string) (model.TrainingRecord, error) {
	if err := auth.CanRecordTraining(actor, employeeID); err != nil {
		return model.TrainingRecord{}, err
	}
	if score < 0 || score > 100 {
		return model.TrainingRecord{}, errors.New("score out of range")
	}
	r := model.TrainingRecord{ID: fmt.Sprintf("training-%s-%s", employeeID, policyID), EmployeeID: employeeID, PolicyID: policyID, CompletedAt: s.now(), Score: score, Trainer: trainer, Notes: notes}
	if err := model.ValidateTrainingRecord(r); err != nil {
		return r, err
	}
	return r, s.store.SaveTraining(r)
}

func (s *Service) Get(actor model.Employee, id string) (model.TrainingRecord, error) {
	r, err := s.store.GetTraining(id)
	if err != nil {
		return r, err
	}
	if err := auth.CanRecordTraining(actor, r.EmployeeID); err != nil {
		return r, err
	}
	return r, nil
}

func (s *Service) ForEmployee(actor model.Employee, employeeID string) ([]model.TrainingRecord, error) {
	if err := auth.CanRecordTraining(actor, employeeID); err != nil {
		return nil, err
	}
	all, err := s.store.ListTraining()
	if err != nil {
		return nil, err
	}
	result := make([]model.TrainingRecord, 0)
	for _, record := range all {
		if record.EmployeeID == employeeID {
			result = append(result, record)
		}
	}
	return result, nil
}

func (s *Service) CompletionRate(employeeID string, policyIDs []string) (float64, error) {
	if len(policyIDs) == 0 {
		return 0, errors.New("policy list cannot be empty")
	}
	completed, err := s.store.ListTraining()
	if err != nil {
		return 0, err
	}
	done := 0
	for _, policyID := range policyIDs {
		for _, record := range completed {
			if record.EmployeeID == employeeID && record.PolicyID == policyID && record.Score >= 60 {
				done++
				break
			}
		}
	}
	return float64(done) / float64(len(policyIDs)), nil
}

func (s *Service) SetClock(now func() time.Time) error {
	if now == nil {
		return errors.New("clock cannot be nil")
	}
	s.now = now
	return nil
}
