package storage

import (
	"hospitaldesk/model"
	"sort"
	"strings"
)

func (s *Store) SearchPolicies(filter model.SearchFilter) ([]model.PolicyDocument, error) {
	all, err := s.ListPolicies()
	if err != nil {
		return nil, err
	}
	result := make([]model.PolicyDocument, 0, len(all))
	for _, p := range all {
		if filter.Department != "" && p.Department != filter.Department {
			continue
		}
		if filter.Status != "" && p.Status != filter.Status {
			continue
		}
		if filter.Query != "" && !strings.Contains(strings.ToLower(p.Title+" "+p.Content), strings.ToLower(filter.Query)) {
			continue
		}
		if !filter.IncludeExpired && p.Status == model.PolicyPublished && !p.EffectiveUntil.IsZero() {
			continue
		}
		result = append(result, p)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].UpdatedAt.After(result[j].UpdatedAt) })
	return result, nil
}

func (s *Store) EmployeeDownloads(employeeID string) ([]model.DownloadRecord, error) {
	all, err := s.ListDownloads()
	if err != nil {
		return nil, err
	}
	result := make([]model.DownloadRecord, 0)
	for _, record := range all {
		if record.EmployeeID == employeeID {
			result = append(result, record)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].DownloadedAt.Before(result[j].DownloadedAt) })
	return result, nil
}

func (s *Store) PolicyReviews(policyID string) ([]model.ReviewEvent, error) {
	all, err := s.ListReviews()
	if err != nil {
		return nil, err
	}
	result := make([]model.ReviewEvent, 0)
	for _, event := range all {
		if event.PolicyID == policyID {
			result = append(result, event)
		}
	}
	return result, nil
}
