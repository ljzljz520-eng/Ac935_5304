package storage

import "hospitaldesk/model"

func (s *Store) SaveEmployee(e model.Employee) error {
	return put(s, bucketNames["employees"], e.ID, e)
}
func (s *Store) GetEmployee(id string) (model.Employee, error) {
	return get[model.Employee](s, bucketNames["employees"], id)
}
func (s *Store) ListEmployees() ([]model.Employee, error) {
	return list[model.Employee](s, bucketNames["employees"])
}
func (s *Store) DeleteEmployee(id string) error { return remove(s, bucketNames["employees"], id) }

func (s *Store) SavePolicy(p model.PolicyDocument) error {
	return put(s, bucketNames["policies"], p.ID, p)
}
func (s *Store) GetPolicy(id string) (model.PolicyDocument, error) {
	return get[model.PolicyDocument](s, bucketNames["policies"], id)
}
func (s *Store) ListPolicies() ([]model.PolicyDocument, error) {
	return list[model.PolicyDocument](s, bucketNames["policies"])
}
func (s *Store) DeletePolicy(id string) error { return remove(s, bucketNames["policies"], id) }

func (s *Store) SaveTraining(r model.TrainingRecord) error {
	return put(s, bucketNames["training"], r.ID, r)
}
func (s *Store) GetTraining(id string) (model.TrainingRecord, error) {
	return get[model.TrainingRecord](s, bucketNames["training"], id)
}
func (s *Store) ListTraining() ([]model.TrainingRecord, error) {
	return list[model.TrainingRecord](s, bucketNames["training"])
}

func (s *Store) SaveSchedule(f model.ScheduleFile) error {
	return put(s, bucketNames["schedules"], f.ID, f)
}
func (s *Store) GetSchedule(id string) (model.ScheduleFile, error) {
	return get[model.ScheduleFile](s, bucketNames["schedules"], id)
}
func (s *Store) ListSchedules() ([]model.ScheduleFile, error) {
	return list[model.ScheduleFile](s, bucketNames["schedules"])
}

func (s *Store) SaveDownload(r model.DownloadRecord) error {
	return put(s, bucketNames["downloads"], r.ID, r)
}
func (s *Store) ListDownloads() ([]model.DownloadRecord, error) {
	return list[model.DownloadRecord](s, bucketNames["downloads"])
}

func (s *Store) SaveShare(link model.ShareLink) error {
	return put(s, bucketNames["shares"], link.Token, link)
}
func (s *Store) GetShare(token string) (model.ShareLink, error) {
	return get[model.ShareLink](s, bucketNames["shares"], token)
}

func (s *Store) SaveReview(event model.ReviewEvent) error {
	return put(s, bucketNames["reviews"], event.ID, event)
}
func (s *Store) ListReviews() ([]model.ReviewEvent, error) {
	return list[model.ReviewEvent](s, bucketNames["reviews"])
}
