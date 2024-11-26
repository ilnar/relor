package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/gemlab-dev/relor/internal/model"
	"github.com/google/uuid"
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
	// TODO: harden garbage collection.
	s.cleanup()

	s.mu.Lock()
	defer s.mu.Unlock()
	s.s[j.ID().ID] = j
	return nil
}

func (s *JobStorage) getExpitedJobs() ([]uuid.UUID, error) {
	dt := time.Now()
	ids := []uuid.UUID{}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for id, j := range s.s {
		if j.ExpiresAt().Before(dt) {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func (s *JobStorage) cleanup() {
	ids, err := s.getExpitedJobs()
	if err != nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, id := range ids {
		delete(s.s, id)
	}
}
