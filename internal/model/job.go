package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type JobID struct {
	ID             uuid.UUID
	WorkflowID     uuid.UUID
	WorkflowAction string
}

type Job struct {
	jid       JobID
	createdAt time.Time
	claimedAt time.Time
	closedAt  time.Time
	closeMsg  string
}

func NewJob(jid JobID, createdAt time.Time) *Job {
	return &Job{
		jid:       jid,
		createdAt: createdAt,
	}
}

func (j *Job) ID() JobID {
	return j.jid
}

func (j *Job) ClaimAt(t time.Time) error {
	if !j.claimedAt.IsZero() {
		return fmt.Errorf("job already claimed at %v", j.claimedAt)
	}
	if j.createdAt.After(t) {
		return fmt.Errorf("job created at %v, can't be claimed at %v", j.createdAt, t)
	}
	j.claimedAt = t
	return nil
}

func (j *Job) CloseAt(t time.Time, msg string) error {
	if !j.closedAt.IsZero() {
		return fmt.Errorf("job already closed at %v", j.closedAt)
	}
	if j.claimedAt.IsZero() {
		return fmt.Errorf("job not claimed")
	}
	if j.claimedAt.After(t) {
		return fmt.Errorf("job claimed at %v, can't be closed at %v", j.claimedAt, t)
	}
	j.closedAt = t
	j.closeMsg = msg
	return nil
}
