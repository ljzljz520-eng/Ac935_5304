package model

import (
	"fmt"
	"sort"
	"strings"
)

func NormalizeEmployee(e Employee) Employee {
	e.ID = strings.TrimSpace(e.ID)
	e.Name = strings.TrimSpace(e.Name)
	e.Department = strings.ToLower(strings.TrimSpace(e.Department))
	return e
}

func NormalizePolicy(p PolicyDocument) PolicyDocument {
	p.ID = strings.TrimSpace(p.ID)
	p.Title = strings.TrimSpace(p.Title)
	p.Department = strings.ToLower(strings.TrimSpace(p.Department))
	p.Content = strings.TrimSpace(p.Content)
	return p
}

func PolicyLabel(p PolicyDocument) string {
	return fmt.Sprintf("%s v%d (%s)", p.Title, p.Version, p.Status)
}

func PolicySnippet(p PolicyDocument, limit int) string {
	if limit < 0 {
		limit = 0
	}
	content := strings.TrimSpace(p.Content)
	if len(content) <= limit {
		return content
	}
	if limit < 4 {
		return content[:limit]
	}
	return content[:limit-3] + "..."
}

func StatusRank(status PolicyStatus) int {
	switch status {
	case PolicyDraft:
		return 1
	case PolicyPending:
		return 2
	case PolicyPublished:
		return 3
	case PolicyArchived:
		return 4
	default:
		return 0
	}
}

func SortPoliciesByStatus(policies []PolicyDocument) []PolicyDocument {
	result := append([]PolicyDocument(nil), policies...)
	sort.Slice(result, func(i, j int) bool {
		left, right := StatusRank(result[i].Status), StatusRank(result[j].Status)
		if left == right {
			return result[i].Title < result[j].Title
		}
		return left < right
	})
	return result
}

func DistinctDepartments(employees []Employee) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)
	for _, employee := range employees {
		if !seen[employee.Department] {
			seen[employee.Department] = true
			result = append(result, employee.Department)
		}
	}
	sort.Strings(result)
	return result
}

func DepartmentRoster(employees []Employee, department string) []Employee {
	result := make([]Employee, 0)
	for _, employee := range employees {
		if strings.EqualFold(employee.Department, department) && employee.Active {
			result = append(result, employee)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}

func TrainingGrade(score int) string {
	if score >= 90 {
		return "excellent"
	}
	if score >= 80 {
		return "good"
	}
	if score >= 60 {
		return "pass"
	}
	return "retry"
}

func IsRoleAtLeast(role UserRole, required UserRole) bool {
	rank := map[UserRole]int{RoleEmployee: 1, RoleSupervisor: 2, RoleDirector: 3}
	return rank[role] >= rank[required]
}

func DownloadDescription(record DownloadRecord) string {
	status := "denied"
	if record.Allowed {
		status = "allowed"
	}
	return fmt.Sprintf("%s %s %s", record.ResourceType, record.ResourceID, status)
}
