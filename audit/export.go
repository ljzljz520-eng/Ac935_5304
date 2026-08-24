package audit

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"hospitaldesk/model"
	"sort"
	"strings"
	"time"
)

func CSV(records []model.DownloadRecord) (string, error) {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	if err := writer.Write([]string{"id", "employee_id", "resource_id", "resource_type", "downloaded_at", "allowed", "reason"}); err != nil {
		return "", err
	}
	for _, record := range records {
		if err := writer.Write([]string{record.ID, record.EmployeeID, record.ResourceID, record.ResourceType, record.DownloadedAt.Format(time.RFC3339), fmt.Sprintf("%t", record.Allowed), record.Reason}); err != nil {
			return "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func ByResourceType(records []model.DownloadRecord, resourceType string) []model.DownloadRecord {
	result := make([]model.DownloadRecord, 0)
	for _, record := range records {
		if strings.EqualFold(record.ResourceType, resourceType) {
			result = append(result, record)
		}
	}
	return result
}

func ByDay(records []model.DownloadRecord, day time.Time) []model.DownloadRecord {
	result := make([]model.DownloadRecord, 0)
	for _, record := range records {
		y, m, d := record.DownloadedAt.Date()
		yy, mm, dd := day.Date()
		if y == yy && m == mm && d == dd {
			result = append(result, record)
		}
	}
	return result
}

func SortChronological(records []model.DownloadRecord) []model.DownloadRecord {
	result := append([]model.DownloadRecord(nil), records...)
	sort.Slice(result, func(i, j int) bool { return result[i].DownloadedAt.Before(result[j].DownloadedAt) })
	return result
}

func Reasons(records []model.DownloadRecord) map[string]int {
	result := make(map[string]int)
	for _, record := range records {
		if !record.Allowed {
			result[record.Reason]++
		}
	}
	return result
}
