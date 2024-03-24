package storage

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/ilnar/wf/internal/model"
)

type JobStorage struct {
	s  map[uuid.UUID]*model.Job
	mu sync.RWMutex
}

func NewJobStorage() *JobStorage {
	return &JobStorage{s: make(map[uuid.UUID]*model.Job)}
}

func (s *JobStorage) Get(id uuid.UUID) (*model.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.s[id]
	if !ok {
		return nil, fmt.Errorf("job not found")
	}
	return j, nil
}

func (s *JobStorage) Save(j *model.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s[j.ID().ID] = j
	return nil
}
