package service

import (
	"encoding/json"
	"errors"
	"hospitaldesk/audit"
	"hospitaldesk/model"
	"sort"
	"strings"
	"time"
)

type PolicySummary struct {
	Total        int            `json:"total"`
	Draft        int            `json:"draft"`
	Pending      int            `json:"pending"`
	Published    int            `json:"published"`
	Expired      int            `json:"expired"`
	ByDepartment map[string]int `json:"by_department"`
}

type TrainingSummary struct {
	Total      int            `json:"total"`
	Passed     int            `json:"passed"`
	Failed     int            `json:"failed"`
	Average    float64        `json:"average"`
	ByEmployee map[string]int `json:"by_employee"`
}

type ScheduleSummary struct {
	Total       int      `json:"total"`
	Published   int      `json:"published"`
	Periods     []string `json:"periods"`
	Departments []string `json:"departments"`
}

type AccessSnapshot struct {
	GeneratedAt time.Time       `json:"generated_at"`
	EmployeeID  string          `json:"employee_id"`
	Policies    PolicySummary   `json:"policies"`
	Training    TrainingSummary `json:"training"`
	Schedules   ScheduleSummary `json:"schedules"`
	Downloads   audit.Summary   `json:"downloads"`
}

func SummarizePolicies(policies []model.PolicyDocument, now time.Time) PolicySummary {
	summary := PolicySummary{ByDepartment: make(map[string]int)}
	for _, p := range policies {
		summary.Total++
		summary.ByDepartment[p.Department]++
		switch p.Status {
		case model.PolicyDraft:
			summary.Draft++
		case model.PolicyPending:
			summary.Pending++
		case model.PolicyPublished:
			summary.Published++
			if !p.EffectiveUntil.IsZero() && now.After(p.EffectiveUntil) {
				summary.Expired++
			}
		case model.PolicyArchived:
		}
	}
	return summary
}

func SummarizeTraining(records []model.TrainingRecord) TrainingSummary {
	summary := TrainingSummary{ByEmployee: make(map[string]int)}
	totalScore := 0
	for _, record := range records {
		summary.Total++
		summary.ByEmployee[record.EmployeeID]++
		totalScore += record.Score
		if record.Score >= 60 {
			summary.Passed++
		} else {
			summary.Failed++
		}
	}
	if summary.Total > 0 {
		summary.Average = float64(totalScore) / float64(summary.Total)
	}
	return summary
}

func SummarizeSchedules(files []model.ScheduleFile) ScheduleSummary {
	summary := ScheduleSummary{Periods: make([]string, 0), Departments: make([]string, 0)}
	periods := make(map[string]bool)
	departments := make(map[string]bool)
	for _, file := range files {
		summary.Total++
		if file.Status == model.PolicyPublished {
			summary.Published++
		}
		if !periods[file.Period] {
			periods[file.Period] = true
			summary.Periods = append(summary.Periods, file.Period)
		}
		if !departments[file.Department] {
			departments[file.Department] = true
			summary.Departments = append(summary.Departments, file.Department)
		}
	}
	sort.Strings(summary.Periods)
	sort.Strings(summary.Departments)
	return summary
}

func (d *Desk) BuildSnapshot(actor model.Employee, now time.Time) (AccessSnapshot, error) {
	if !actor.Active {
		return AccessSnapshot{}, errors.New("inactive employee")
	}
	policies, err := d.Store.ListPolicies()
	if err != nil {
		return AccessSnapshot{}, err
	}
	training, err := d.Store.ListTraining()
	if err != nil {
		return AccessSnapshot{}, err
	}
	schedules, err := d.Store.ListSchedules()
	if err != nil {
		return AccessSnapshot{}, err
	}
	downloads, err := d.Audit.History(actor.ID)
	if err != nil {
		return AccessSnapshot{}, err
	}
	return AccessSnapshot{GeneratedAt: now, EmployeeID: actor.ID, Policies: SummarizePolicies(policies, now), Training: SummarizeTraining(training), Schedules: SummarizeSchedules(schedules), Downloads: audit.BuildSummary(downloads)}, nil
}

func EncodeSnapshot(snapshot AccessSnapshot) ([]byte, error) {
	return json.MarshalIndent(snapshot, "", "  ")
}

func FilterDepartments(summary PolicySummary, allowed []string) PolicySummary {
	if len(allowed) == 0 {
		return summary
	}
	selected := make(map[string]bool)
	for _, name := range allowed {
		selected[strings.TrimSpace(name)] = true
	}
	filtered := summary
	filtered.ByDepartment = make(map[string]int)
	filtered.Total = 0
	for department, count := range summary.ByDepartment {
		if selected[department] {
			filtered.ByDepartment[department] = count
			filtered.Total += count
		}
	}
	return filtered
}

func HighestTrainingScore(records []model.TrainingRecord) (model.TrainingRecord, bool) {
	if len(records) == 0 {
		return model.TrainingRecord{}, false
	}
	highest := records[0]
	for _, record := range records[1:] {
		if record.Score > highest.Score {
			highest = record
		}
	}
	return highest, true
}

func LatestPolicy(policies []model.PolicyDocument) (model.PolicyDocument, bool) {
	if len(policies) == 0 {
		return model.PolicyDocument{}, false
	}
	latest := policies[0]
	for _, p := range policies[1:] {
		if p.UpdatedAt.After(latest.UpdatedAt) {
			latest = p
		}
	}
	return latest, true
}

func VisiblePolicyCount(actor model.Employee, policies []model.PolicyDocument, now time.Time) int {
	count := 0
	for _, p := range policies {
		if actor.CanViewDepartment(p.Department) && p.IsCurrentlyValid(now) {
			count++
		}
	}
	return count
}
