package port

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Upsert(ctx context.Context, p Port) error {
	if err := p.Validate(); err != nil {
		return err
	}
	return s.repo.Upsert(ctx, p)
}

func (s *Service) Get(ctx context.Context, id string) (Port, bool) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Count() int {
	return s.repo.Count()
}
