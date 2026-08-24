package audit

import (
	"hospitaldesk/model"
	"sort"
)

type Summary struct {
	Total   int
	Allowed int
	Denied  int
	ByType  map[string]int
}

func BuildSummary(records []model.DownloadRecord) Summary {
	s := Summary{ByType: make(map[string]int)}
	for _, record := range records {
		s.Total++
		s.ByType[record.ResourceType]++
		if record.Allowed {
			s.Allowed++
		} else {
			s.Denied++
		}
	}
	return s
}

func Recent(records []model.DownloadRecord, limit int) []model.DownloadRecord {
	copyRecords := append([]model.DownloadRecord(nil), records...)
	sort.Slice(copyRecords, func(i, j int) bool { return copyRecords[i].DownloadedAt.After(copyRecords[j].DownloadedAt) })
	if limit < 0 {
		limit = 0
	}
	if limit > len(copyRecords) {
		limit = len(copyRecords)
	}
	return copyRecords[:limit]
}

func AllowedOnly(records []model.DownloadRecord) []model.DownloadRecord {
	result := make([]model.DownloadRecord, 0)
	for _, r := range records {
		if r.Allowed {
			result = append(result, r)
		}
	}
	return result
}
