package model

import "time"

type UserRole string

const (
	RoleDirector   UserRole = "director"
	RoleSupervisor UserRole = "supervisor"
	RoleEmployee   UserRole = "employee"
)

type PolicyStatus string

const (
	PolicyDraft     PolicyStatus = "draft"
	PolicyPending   PolicyStatus = "pending_review"
	PolicyPublished PolicyStatus = "published"
	PolicyArchived  PolicyStatus = "archived"
)

type Employee struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Department string    `json:"department"`
	Role       UserRole  `json:"role"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
}

type PolicyDocument struct {
	ID             string       `json:"id"`
	Title          string       `json:"title"`
	Department     string       `json:"department"`
	Content        string       `json:"content"`
	Status         PolicyStatus `json:"status"`
	Version        int          `json:"version"`
	EffectiveFrom  time.Time    `json:"effective_from"`
	EffectiveUntil time.Time    `json:"effective_until"`
	CreatedBy      string       `json:"created_by"`
	ApprovedBy     string       `json:"approved_by"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type TrainingRecord struct {
	ID          string    `json:"id"`
	EmployeeID  string    `json:"employee_id"`
	PolicyID    string    `json:"policy_id"`
	CompletedAt time.Time `json:"completed_at"`
	Score       int       `json:"score"`
	Trainer     string    `json:"trainer"`
	Notes       string    `json:"notes"`
}

type ScheduleFile struct {
	ID          string       `json:"id"`
	Department  string       `json:"department"`
	Period      string       `json:"period"`
	Content     string       `json:"content"`
	Status      PolicyStatus `json:"status"`
	PublishedAt time.Time    `json:"published_at"`
	CreatedBy   string       `json:"created_by"`
}

type DownloadRecord struct {
	ID           string    `json:"id"`
	EmployeeID   string    `json:"employee_id"`
	ResourceID   string    `json:"resource_id"`
	ResourceType string    `json:"resource_type"`
	DownloadedAt time.Time `json:"downloaded_at"`
	Allowed      bool      `json:"allowed"`
	Reason       string    `json:"reason"`
}

type ShareLink struct {
	Token        string    `json:"token"`
	ResourceID   string    `json:"resource_id"`
	ResourceType string    `json:"resource_type"`
	Department   string    `json:"department"`
	CreatedBy    string    `json:"created_by"`
	ExpiresAt    time.Time `json:"expires_at"`
	Revoked      bool      `json:"revoked"`
}

type ReviewEvent struct {
	ID       string    `json:"id"`
	PolicyID string    `json:"policy_id"`
	ActorID  string    `json:"actor_id"`
	Action   string    `json:"action"`
	At       time.Time `json:"at"`
	Comment  string    `json:"comment"`
}

type SearchFilter struct {
	Department     string
	Status         PolicyStatus
	Query          string
	IncludeExpired bool
}

func (p PolicyDocument) IsCurrentlyValid(now time.Time) bool {
	if p.Status != PolicyPublished {
		return false
	}
	if !p.EffectiveFrom.IsZero() && now.Before(p.EffectiveFrom) {
		return false
	}
	if !p.EffectiveUntil.IsZero() && now.After(p.EffectiveUntil) {
		return false
	}
	return true
}

func (e Employee) CanManagePolicies() bool {
	return e.Active && (e.Role == RoleDirector || e.Role == RoleSupervisor)
}

func (e Employee) CanViewDepartment(dept string) bool {
	if !e.Active {
		return false
	}
	if e.Role == RoleDirector {
		return true
	}
	return e.Department == dept
}

func (s ScheduleFile) IsPublished() bool { return s.Status == PolicyPublished }

func (r DownloadRecord) IsSuccessful() bool { return r.Allowed && r.Reason == "" }
