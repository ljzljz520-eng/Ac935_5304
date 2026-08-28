package service

import (
	"errors"
	"hospitaldesk/audit"
	"hospitaldesk/model"
	"hospitaldesk/policy"
	"hospitaldesk/schedule"
	"hospitaldesk/storage"
	"hospitaldesk/training"
	"time"
)

type Desk struct {
	Store     *storage.Store
	Policies  *policy.Manager
	Shares    *policy.ShareService
	Training  *training.Service
	Schedules *schedule.Service
	Audit     *audit.Logger
}

func New(store *storage.Store) *Desk {
	return &Desk{Store: store, Policies: policy.NewManager(store), Shares: policy.NewShareService(store), Training: training.NewService(store), Schedules: schedule.NewService(store), Audit: audit.NewLogger(store)}
}

func (d *Desk) RegisterEmployee(employee model.Employee) error {
	if err := model.ValidateEmployee(employee); err != nil {
		return err
	}
	if employee.CreatedAt.IsZero() {
		employee.CreatedAt = time.Date(2026, 1, 2, 9, 0, 0, 0, time.UTC)
	}
	return d.Store.SaveEmployee(employee)
}

func (d *Desk) Employee(id string) (model.Employee, error) { return d.Store.GetEmployee(id) }

func (d *Desk) PublishPolicyWorkflow(actor model.Employee, title, dept, content string, from, until time.Time) (model.PolicyDocument, error) {
	p, err := d.Policies.CreateDraft(actor, title, dept, content)
	if err != nil {
		return p, err
	}
	p, err = d.Policies.SubmitForReview(actor, p.ID)
	if err != nil {
		return p, err
	}
	return d.Policies.Approve(actor, p.ID, from, until, "director review complete")
}

func (d *Desk) RecordPolicyDownload(actor model.Employee, policyID string) (model.DownloadRecord, error) {
	p, err := d.Policies.GetForEmployee(actor, policyID)
	if err != nil {
		record, _ := d.Audit.Record(actor.ID, policyID, "policy", false, err.Error())
		return record, err
	}
	return d.Audit.Record(actor.ID, p.ID, "policy", true, "")
}

func (d *Desk) RecordShareDownload(actor model.Employee, token string) (model.DownloadRecord, error) {
	p, err := d.Shares.Download(actor, token)
	if err != nil {
		resourceID := token
		if link, linkErr := d.Store.GetShare(token); linkErr == nil {
			resourceID = link.ResourceID
		}
		record, _ := d.Audit.Record(actor.ID, resourceID, "policy", false, err.Error())
		return record, err
	}
	return d.Audit.Record(actor.ID, p.ID, "policy", true, "")
}

func (d *Desk) CompleteTraining(actor model.Employee, employeeID, policyID string, score int, trainer string) (model.TrainingRecord, error) {
	return d.Training.Record(actor, employeeID, policyID, score, trainer, "completed in department session")
}

func (d *Desk) PublishScheduleWorkflow(actor model.Employee, dept, period, content string) (model.ScheduleFile, error) {
	f, err := d.Schedules.Create(actor, dept, period, content)
	if err != nil {
		return f, err
	}
	return d.Schedules.Publish(actor, f.ID)
}

func (d *Desk) ViewSchedule(actor model.Employee, id string) (model.ScheduleFile, error) {
	return d.Schedules.View(actor, id)
}

func (d *Desk) Dashboard(actor model.Employee) (Dashboard, error) {
	if !actor.Active {
		return Dashboard{}, errors.New("inactive employee")
	}
	policies, err := d.Policies.Search(actor, model.SearchFilter{Status: model.PolicyPublished})
	if err != nil {
		return Dashboard{}, err
	}
	schedules, err := d.Schedules.Search(actor, actor.Department, "")
	if err != nil {
		return Dashboard{}, err
	}
	trainingRecords, err := d.Training.ForEmployee(actor, actor.ID)
	if err != nil {
		return Dashboard{}, err
	}
	downloads, err := d.Audit.History(actor.ID)
	if err != nil {
		return Dashboard{}, err
	}
	return Dashboard{Policies: policies, Schedules: schedules, Training: trainingRecords, Downloads: downloads, DownloadSummary: audit.BuildSummary(downloads)}, nil
}

type Dashboard struct {
	Policies        []model.PolicyDocument
	Schedules       []model.ScheduleFile
	Training        []model.TrainingRecord
	Downloads       []model.DownloadRecord
	DownloadSummary audit.Summary
}
