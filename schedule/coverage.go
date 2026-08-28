package schedule

import (
	"errors"
	"hospitaldesk/model"
	"sort"
	"strconv"
	"strings"
)

type Shift struct {
	Day        string
	Start      string
	End        string
	EmployeeID string
}

type CoverageReport struct {
	Period      string
	Department  string
	Shifts      int
	Employees   map[string]int
	MissingDays []string
}

func ParseShifts(content string) ([]Shift, error) {
	lines := strings.Split(content, "\n")
	shifts := make([]Shift, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 4 {
			return nil, errors.New("shift line must have day,start,end,employee")
		}
		if parts[0] == "" || parts[1] == "" || parts[2] == "" || parts[3] == "" {
			return nil, errors.New("shift fields cannot be empty")
		}
		shifts = append(shifts, Shift{Day: strings.TrimSpace(parts[0]), Start: strings.TrimSpace(parts[1]), End: strings.TrimSpace(parts[2]), EmployeeID: strings.TrimSpace(parts[3])})
	}
	return shifts, nil
}

func ValidateShiftHours(shift Shift) error {
	start, err := strconv.Atoi(strings.ReplaceAll(shift.Start, ":", ""))
	if err != nil {
		return err
	}
	end, err := strconv.Atoi(strings.ReplaceAll(shift.End, ":", ""))
	if err != nil {
		return err
	}
	if start < 0 || end < 0 || start > 2359 || end > 2359 {
		return errors.New("shift time is invalid")
	}
	if end <= start {
		return errors.New("shift must end after start")
	}
	return nil
}

func BuildCoverage(file model.ScheduleFile, requiredDays []string) (CoverageReport, error) {
	shifts, err := ParseShifts(file.Content)
	if err != nil {
		return CoverageReport{}, err
	}
	report := CoverageReport{Period: file.Period, Department: file.Department, Shifts: len(shifts), Employees: make(map[string]int), MissingDays: make([]string, 0)}
	days := make(map[string]bool)
	for _, shift := range shifts {
		if err := ValidateShiftHours(shift); err != nil {
			return report, err
		}
		days[shift.Day] = true
		report.Employees[shift.EmployeeID]++
	}
	for _, day := range requiredDays {
		if !days[day] {
			report.MissingDays = append(report.MissingDays, day)
		}
	}
	return report, nil
}

func SortShifts(shifts []Shift) []Shift {
	result := append([]Shift(nil), shifts...)
	sort.Slice(result, func(i, j int) bool {
		if result[i].Day == result[j].Day {
			return result[i].Start < result[j].Start
		}
		return result[i].Day < result[j].Day
	})
	return result
}

func RequiredCoverage(report CoverageReport) bool {
	return report.Shifts > 0 && len(report.MissingDays) == 0
}

func DistinctEmployees(report CoverageReport) int { return len(report.Employees) }

func ScheduleLabel(file model.ScheduleFile) string {
	return strings.TrimSpace(file.Department) + " / " + strings.TrimSpace(file.Period)
}
