package memory

import (
	"context"
	"sync"

	"lessonsgo/testmicroservice/internal/domain/port"
)

type Store struct {
	mu    sync.RWMutex
	ports map[string]port.Port
}

func New() *Store {
	return &Store{
		ports: make(map[string]port.Port),
	}
}

func (s *Store) Upsert(ctx context.Context, p port.Port) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.ports[p.ID] = p
	return nil
}

func (s *Store) Get(_ context.Context, id string) (port.Port, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.ports[id]
	return p, ok
}

func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.ports)
}
