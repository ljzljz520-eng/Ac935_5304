package training

import (
	"hospitaldesk/model"
	"sort"
	"strings"
)

func FilterRecords(records []model.TrainingRecord, query string, minimumScore int) []model.TrainingRecord {
	result := make([]model.TrainingRecord, 0, len(records))
	needle := strings.ToLower(strings.TrimSpace(query))
	for _, record := range records {
		if record.Score < minimumScore {
			continue
		}
		if needle != "" && !strings.Contains(strings.ToLower(record.EmployeeID+" "+record.PolicyID+" "+record.Notes), needle) {
			continue
		}
		result = append(result, record)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Score > result[j].Score })
	return result
}

func Summarize(records []model.TrainingRecord) (int, int) {
	passed := 0
	total := 0
	for _, record := range records {
		total++
		if record.Score >= 60 {
			passed++
		}
	}
	return passed, total
}
