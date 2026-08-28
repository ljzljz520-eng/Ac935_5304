package storage

import (
	"errors"
	"fmt"
	"hospitaldesk/model"
	"sort"
	"strings"
)

type Inventory struct {
	Employees int
	Policies  int
	Training  int
	Schedules int
	Downloads int
	Shares    int
	Reviews   int
}

func (s *Store) Inventory() (Inventory, error) {
	values := []struct {
		name   string
		target *int
	}{{"employees", nil}, {"policies", nil}, {"training", nil}, {"schedules", nil}, {"downloads", nil}, {"shares", nil}, {"reviews", nil}}
	counts := make([]int, len(values))
	for i := range values {
		count, err := s.Count(values[i].name)
		if err != nil {
			return Inventory{}, err
		}
		counts[i] = count
	}
	return Inventory{Employees: counts[0], Policies: counts[1], Training: counts[2], Schedules: counts[3], Downloads: counts[4], Shares: counts[5], Reviews: counts[6]}, nil
}

func (s *Store) RemovePolicyCascade(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("policy id is required")
	}
	if err := s.DeletePolicy(id); err != nil {
		return err
	}
	training, err := s.ListTraining()
	if err != nil {
		return err
	}
	for _, record := range training {
		if record.PolicyID == id {
			if err := remove(s, bucketNames["training"], record.ID); err != nil {
				return err
			}
		}
	}
	reviews, err := s.ListReviews()
	if err != nil {
		return err
	}
	for _, review := range reviews {
		if review.PolicyID == id {
			if err := remove(s, bucketNames["reviews"], review.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Store) PolicyTimeline(id string) ([]model.ReviewEvent, error) {
	events, err := s.PolicyReviews(id)
	if err != nil {
		return nil, err
	}
	sort.Slice(events, func(i, j int) bool { return events[i].At.Before(events[j].At) })
	return events, nil
}

func (s *Store) ValidateReferences() error {
	policies, err := s.ListPolicies()
	if err != nil {
		return err
	}
	policyIDs := make(map[string]bool)
	for _, p := range policies {
		policyIDs[p.ID] = true
	}
	training, err := s.ListTraining()
	if err != nil {
		return err
	}
	for _, record := range training {
		if !policyIDs[record.PolicyID] {
			return fmt.Errorf("training %s references missing policy", record.ID)
		}
	}
	schedules, err := s.ListSchedules()
	if err != nil {
		return err
	}
	for _, file := range schedules {
		if strings.TrimSpace(file.Department) == "" {
			return fmt.Errorf("schedule %s has no department", file.ID)
		}
	}
	return nil
}

func (s *Store) PurgeRevokedShares() (int, error) {
	shares, err := list[model.ShareLink](s, bucketNames["shares"])
	if err != nil {
		return 0, err
	}
	removed := 0
	for _, share := range shares {
		if share.Revoked {
			if err := remove(s, bucketNames["shares"], share.Token); err != nil {
				return removed, err
			}
			removed++
		}
	}
	return removed, nil
}

func BucketNames() []string {
	result := make([]string, 0, len(bucketNames))
	for name := range bucketNames {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

func (s *Store) HasPolicy(id string) bool { _, err := s.GetPolicy(id); return err == nil }
