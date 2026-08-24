package audit

import (
	"errors"
	"fmt"
	"hospitaldesk/model"
	"hospitaldesk/storage"
	"strings"
	"time"
)

type Logger struct {
	store *storage.Store
	now   func() time.Time
}

func NewLogger(store *storage.Store) *Logger {
	return &Logger{store: store, now: func() time.Time { return time.Date(2026, 1, 2, 9, 0, 0, 0, time.UTC) }}
}

func (l *Logger) Record(employeeID, resourceID, resourceType string, allowed bool, reason string) (model.DownloadRecord, error) {
	if employeeID == "" || resourceID == "" || resourceType == "" {
		return model.DownloadRecord{}, errors.New("download identity is required")
	}
	if allowed {
		reason = ""
	} else if strings.TrimSpace(reason) == "" {
		reason = "denied"
	}
	r := model.DownloadRecord{ID: fmt.Sprintf("download-%s-%d", employeeID, l.now().UnixNano()), EmployeeID: employeeID, ResourceID: resourceID, ResourceType: resourceType, DownloadedAt: l.now(), Allowed: allowed, Reason: reason}
	return r, l.store.SaveDownload(r)
}

func (l *Logger) History(employeeID string) ([]model.DownloadRecord, error) {
	return l.store.EmployeeDownloads(employeeID)
}

func (l *Logger) CountDenied(records []model.DownloadRecord) int {
	count := 0
	for _, record := range records {
		if !record.Allowed {
			count++
		}
	}
	return count
}

func (l *Logger) SetClock(now func() time.Time) error {
	if now == nil {
		return errors.New("clock cannot be nil")
	}
	l.now = now
	return nil
}
