package memory

import (
	"context"
	"sync"

	"lessonsgo/testpaymentgeteway/internal/domain/payment"
)

type Store struct {
	mu       sync.RWMutex
	payments map[string]payment.Payment
}

func New() *Store {
	return &Store{
		payments: make(map[string]payment.Payment),
	}
}

func (s *Store) Save(ctx context.Context, p payment.Payment) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.payments[p.ID] = p
	return nil
}

func (s *Store) GetByID(_ context.Context, id string) (payment.Payment, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.payments[id]
	return p, ok
}

func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.payments)
}
