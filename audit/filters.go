package audit

import (
	"hospitaldesk/model"
	"sort"
	"strings"
	"time"
)

type Filter struct {
	EmployeeID   string
	ResourceType string
	Allowed      *bool
	From         time.Time
	Until        time.Time
	Reason       string
}

func FilterRecords(records []model.DownloadRecord, filter Filter) []model.DownloadRecord {
	result := make([]model.DownloadRecord, 0)
	for _, record := range records {
		if filter.EmployeeID != "" && record.EmployeeID != filter.EmployeeID {
			continue
		}
		if filter.ResourceType != "" && !strings.EqualFold(record.ResourceType, filter.ResourceType) {
			continue
		}
		if filter.Allowed != nil && record.Allowed != *filter.Allowed {
			continue
		}
		if !filter.From.IsZero() && record.DownloadedAt.Before(filter.From) {
			continue
		}
		if !filter.Until.IsZero() && record.DownloadedAt.After(filter.Until) {
			continue
		}
		if filter.Reason != "" && !strings.Contains(strings.ToLower(record.Reason), strings.ToLower(filter.Reason)) {
			continue
		}
		result = append(result, record)
	}
	return result
}

func GroupByEmployee(records []model.DownloadRecord) map[string][]model.DownloadRecord {
	groups := make(map[string][]model.DownloadRecord)
	for _, record := range records {
		groups[record.EmployeeID] = append(groups[record.EmployeeID], record)
	}
	for employee := range groups {
		groups[employee] = SortChronological(groups[employee])
	}
	return groups
}

func GroupByResource(records []model.DownloadRecord) map[string][]model.DownloadRecord {
	groups := make(map[string][]model.DownloadRecord)
	for _, record := range records {
		groups[record.ResourceID] = append(groups[record.ResourceID], record)
	}
	return groups
}

func DeniedReasons(records []model.DownloadRecord) []string {
	counts := Reasons(records)
	result := make([]string, 0, len(counts))
	for reason := range counts {
		result = append(result, reason)
	}
	sort.Strings(result)
	return result
}

func FirstAccess(records []model.DownloadRecord, resourceID string) (model.DownloadRecord, bool) {
	var found model.DownloadRecord
	has := false
	for _, record := range records {
		if record.ResourceID == resourceID && (!has || record.DownloadedAt.Before(found.DownloadedAt)) {
			found = record
			has = true
		}
	}
	return found, has
}

func LastAccess(records []model.DownloadRecord, resourceID string) (model.DownloadRecord, bool) {
	var found model.DownloadRecord
	has := false
	for _, record := range records {
		if record.ResourceID == resourceID && (!has || record.DownloadedAt.After(found.DownloadedAt)) {
			found = record
			has = true
		}
	}
	return found, has
}

func IsRepeatedDenied(records []model.DownloadRecord, employeeID, resourceID string, threshold int) bool {
	if threshold <= 0 {
		return false
	}
	count := 0
	for _, record := range records {
		if record.EmployeeID == employeeID && record.ResourceID == resourceID && !record.Allowed {
			count++
		}
	}
	return count >= threshold
}

func AllowedRate(records []model.DownloadRecord) float64 {
	if len(records) == 0 {
		return 0
	}
	allowed := 0
	for _, record := range records {
		if record.Allowed {
			allowed++
		}
	}
	return float64(allowed) / float64(len(records))
}

func MergeFilters(filters ...Filter) Filter {
	merged := Filter{}
	for _, filter := range filters {
		if filter.EmployeeID != "" {
			merged.EmployeeID = filter.EmployeeID
		}
		if filter.ResourceType != "" {
			merged.ResourceType = filter.ResourceType
		}
		if filter.Allowed != nil {
			value := *filter.Allowed
			merged.Allowed = &value
		}
		if !filter.From.IsZero() {
			merged.From = filter.From
		}
		if !filter.Until.IsZero() {
			merged.Until = filter.Until
		}
		if filter.Reason != "" {
			merged.Reason = filter.Reason
		}
	}
	return merged
}

func PartitionByAccess(records []model.DownloadRecord) (allowed, denied []model.DownloadRecord) {
	allowed = make([]model.DownloadRecord, 0)
	denied = make([]model.DownloadRecord, 0)
	for _, record := range records {
		if record.Allowed {
			allowed = append(allowed, record)
		} else {
			denied = append(denied, record)
		}
	}
	return allowed, denied
}
