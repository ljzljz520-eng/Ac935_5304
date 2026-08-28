package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"hospitaldesk/model"
	"hospitaldesk/service"
	"hospitaldesk/storage"
	"os"
	"time"
)

func main() {
	path := flag.String("db", "clinicdesk.db", "database path")
	flag.Parse()
	store, err := storage.Open(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer store.Close()
	desk := service.New(store)
	if err := runDemo(desk); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := printStatus(desk); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runDemo(desk *service.Desk) error {
	director := model.Employee{ID: "director-1", Name: "Dr. Lin", Department: "cardiology", Role: model.RoleDirector, Active: true}
	if err := desk.RegisterEmployee(director); err != nil && !isExisting(err) {
		return err
	}
	if _, err := desk.PublishPolicyWorkflow(director, "Medication Safety", "cardiology", "Verify high-risk medications before administration.", time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)); err != nil {
		return err
	}
	if _, err := desk.PublishScheduleWorkflow(director, "cardiology", "2026-W01", "Mon: ward A; Tue: ward B"); err != nil && !isExisting(err) {
		return err
	}
	return nil
}

func printStatus(desk *service.Desk) error {
	employees, err := desk.Store.ListEmployees()
	if err != nil {
		return err
	}
	policies, err := desk.Store.ListPolicies()
	if err != nil {
		return err
	}
	schedules, err := desk.Store.ListSchedules()
	if err != nil {
		return err
	}
	status := map[string]int{"employees": len(employees), "policies": len(policies), "schedules": len(schedules)}
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func isExisting(err error) bool {
	return err != nil && len(err.Error()) > 7 && err.Error()[len(err.Error())-7:] == "exists"
}
