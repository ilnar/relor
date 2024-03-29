package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type labelSet map[string]struct{}

func (l labelSet) Slice() []string {
	s := make([]string, 0, len(l))
	for k := range l {
		s = append(s, k)
	}
	return s
}

type JobID struct {
	ID             uuid.UUID
	WorkflowID     uuid.UUID
	WorkflowAction string
}

type Job struct {
	jid         JobID
	labels      labelSet
	createdAt   time.Time
	claimedAt   time.Time
	closedAt    time.Time
	resultLabel string
}

func NewJob(jid JobID, labels []string, createdAt time.Time) *Job {
	l := make(labelSet, len(labels))
	for _, label := range labels {
		l[label] = struct{}{}
	}
	return &Job{
		jid:       jid,
		labels:    l,
		createdAt: createdAt,
	}
}

func (j *Job) ID() JobID {
	return j.jid
}

func (j *Job) Labels() labelSet {
	return j.labels
}

func (j *Job) ResultLabel() string {
	return j.resultLabel
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

func (j *Job) CloseAt(t time.Time, resultLabel string) error {
	if resultLabel == "" {
		return fmt.Errorf("result label is empty")
	}
	if _, found := j.labels[resultLabel]; !found {
		return fmt.Errorf("result label %q is not among job labels %v", resultLabel, j.labels.Slice())
	}
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
	j.resultLabel = resultLabel
	return nil
}
