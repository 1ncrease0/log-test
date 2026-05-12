package domain

import "time"

type LogStatus string

const (
	LogStatusPending    LogStatus = "pending"
	LogStatusProcessing LogStatus = "processing"
	LogStatusDone       LogStatus = "done"
	LogStatusFailed     LogStatus = "failed"
)

type Log struct {
	ID         int64
	Path       string
	Status     LogStatus
	NodeCount  int
	PortCount  int
	UploadedAt time.Time
}
