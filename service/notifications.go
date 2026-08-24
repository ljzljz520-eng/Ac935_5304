package service

import (
	"errors"
	"hospitaldesk/model"
	"sort"
	"strings"
)

type NoticeKind string

const (
	NoticeReview   NoticeKind = "review"
	NoticeTraining NoticeKind = "training"
	NoticeSchedule NoticeKind = "schedule"
	NoticeAccess   NoticeKind = "access"
)

type Notice struct {
	ID          string     `json:"id"`
	RecipientID string     `json:"recipient_id"`
	Kind        NoticeKind `json:"kind"`
	Subject     string     `json:"subject"`
	Body        string     `json:"body"`
	Read        bool       `json:"read"`
	Sequence    int        `json:"sequence"`
}

type NoticeBook struct {
	notices map[string][]Notice
	next    int
}

func NewNoticeBook() *NoticeBook { return &NoticeBook{notices: make(map[string][]Notice)} }

func (b *NoticeBook) Add(recipient string, kind NoticeKind, subject, body string) (Notice, error) {
	if strings.TrimSpace(recipient) == "" {
		return Notice{}, errors.New("notice recipient is required")
	}
	if subject == "" || body == "" {
		return Notice{}, errors.New("notice content is required")
	}
	if kind != NoticeReview && kind != NoticeTraining && kind != NoticeSchedule && kind != NoticeAccess {
		return Notice{}, errors.New("unknown notice kind")
	}
	b.next++
	notice := Notice{ID: strings.Join([]string{"notice", recipient, string(kind), string(rune(b.next))}, "-"), RecipientID: recipient, Kind: kind, Subject: subject, Body: body, Sequence: b.next}
	b.notices[recipient] = append(b.notices[recipient], notice)
	return notice, nil
}

func (b *NoticeBook) ForEmployee(employeeID string, unreadOnly bool) []Notice {
	items := append([]Notice(nil), b.notices[employeeID]...)
	if unreadOnly {
		filtered := make([]Notice, 0, len(items))
		for _, item := range items {
			if !item.Read {
				filtered = append(filtered, item)
			}
		}
		items = filtered
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Sequence > items[j].Sequence })
	return items
}

func (b *NoticeBook) MarkRead(employeeID, noticeID string) error {
	items, ok := b.notices[employeeID]
	if !ok {
		return errors.New("employee has no notices")
	}
	for index := range items {
		if items[index].ID == noticeID {
			items[index].Read = true
			b.notices[employeeID] = items
			return nil
		}
	}
	return errors.New("notice not found")
}

func (b *NoticeBook) UnreadCount(employeeID string) int {
	count := 0
	for _, item := range b.notices[employeeID] {
		if !item.Read {
			count++
		}
	}
	return count
}

func (b *NoticeBook) ClearEmployee(employeeID string) { delete(b.notices, employeeID) }

func (d *Desk) NotifyPolicyReview(book *NoticeBook, employees []model.Employee, policy model.PolicyDocument) int {
	sent := 0
	for _, employee := range employees {
		if employee.Active && employee.CanViewDepartment(policy.Department) {
			if _, err := book.Add(employee.ID, NoticeReview, "制度审核状态", policy.Title); err == nil {
				sent++
			}
		}
	}
	return sent
}

func (d *Desk) NotifyTrainingResult(book *NoticeBook, employee model.Employee, record model.TrainingRecord) (Notice, error) {
	subject := "培训记录已保存"
	if record.Score < 60 {
		subject = "培训需要补考"
	}
	return book.Add(employee.ID, NoticeTraining, subject, model.TrainingGrade(record.Score))
}

func (d *Desk) NotifySchedule(book *NoticeBook, employee model.Employee, file model.ScheduleFile) (Notice, error) {
	return book.Add(employee.ID, NoticeSchedule, "新排班已发布", file.Period)
}

func NoticeKinds(items []Notice) map[NoticeKind]int {
	counts := make(map[NoticeKind]int)
	for _, item := range items {
		counts[item.Kind]++
	}
	return counts
}

func MergeNotices(books ...*NoticeBook) []Notice {
	result := make([]Notice, 0)
	for _, book := range books {
		for recipient := range book.notices {
			result = append(result, book.notices[recipient]...)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Sequence < result[j].Sequence })
	return result
}
