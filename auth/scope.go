package auth

import (
	"errors"
	"hospitaldesk/model"
	"strings"
)

type Scope struct {
	EmployeeID string
	Department string
	CanRead    bool
	CanWrite   bool
	CanApprove bool
	CanExport  bool
}

func BuildScope(employee model.Employee) Scope {
	scope := Scope{EmployeeID: employee.ID, Department: employee.Department, CanRead: employee.Active}
	if !employee.Active {
		return scope
	}
	if employee.Role == model.RoleDirector {
		scope.CanWrite = true
		scope.CanApprove = true
		scope.CanExport = true
	} else if employee.Role == model.RoleSupervisor {
		scope.CanWrite = true
		scope.CanExport = true
	}
	return scope
}

func RequireRead(scope Scope) error {
	if !scope.CanRead {
		return ErrInactiveEmployee
	}
	return nil
}

func RequireWrite(scope Scope, department string) error {
	if err := RequireRead(scope); err != nil {
		return err
	}
	if !scope.CanWrite {
		return errors.New("write permission required")
	}
	if !scope.CanApprove && !strings.EqualFold(scope.Department, department) {
		return ErrDepartmentDenied
	}
	return nil
}

func RequireApprove(scope Scope) error {
	if err := RequireRead(scope); err != nil {
		return err
	}
	if !scope.CanApprove {
		return errors.New("approval permission required")
	}
	return nil
}

func RequireExport(scope Scope) error {
	if err := RequireRead(scope); err != nil {
		return err
	}
	if !scope.CanExport {
		return errors.New("export permission required")
	}
	return nil
}

func ScopeForDepartment(employee model.Employee, department string) (Scope, error) {
	scope := BuildScope(employee)
	if err := RequireRead(scope); err != nil {
		return scope, err
	}
	if employee.Role != model.RoleDirector && !strings.EqualFold(employee.Department, department) {
		return scope, ErrDepartmentDenied
	}
	return scope, nil
}

func CanShare(scope Scope, department string) bool {
	return scope.CanRead && (scope.CanApprove || strings.EqualFold(scope.Department, department))
}

func PermissionNames(scope Scope) []string {
	permissions := make([]string, 0, 4)
	if scope.CanRead {
		permissions = append(permissions, "read")
	}
	if scope.CanWrite {
		permissions = append(permissions, "write")
	}
	if scope.CanApprove {
		permissions = append(permissions, "approve")
	}
	if scope.CanExport {
		permissions = append(permissions, "export")
	}
	return permissions
}
