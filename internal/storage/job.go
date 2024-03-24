package storage

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ilnar/wf/internal/model"
)

type JobStorage struct {
	s map[uuid.UUID]*model.Job
}

func NewJobStorage() *JobStorage {
	return &JobStorage{s: make(map[uuid.UUID]*model.Job)}
}

func (s *JobStorage) Get(id uuid.UUID) (*model.Job, error) {
	j, ok := s.s[id]
	if !ok {
		return nil, fmt.Errorf("job not found")
	}
	return j, nil
}

func (s *JobStorage) Save(j *model.Job) error {
	s.s[j.ID().ID] = j
	return nil
}
