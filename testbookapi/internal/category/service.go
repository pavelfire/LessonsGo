package category

import (
	"context"
	"strings"

	"lessonsgo/testbookapi/internal/domain"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

type CreateRequest struct {
	Name string `json:"name"`
}

type UpdateRequest struct {
	Name string `json:"name"`
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (domain.Category, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return domain.Category{}, domain.ErrInvalidInput
	}
	return s.repo.CreateCategory(ctx, name)
}

func (s *Service) Get(ctx context.Context, id int64) (domain.Category, error) {
	if id <= 0 {
		return domain.Category{}, domain.ErrInvalidInput
	}
	return s.repo.GetCategory(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]domain.Category, error) {
	return s.repo.ListCategories(ctx)
}

func (s *Service) Update(ctx context.Context, id int64, req UpdateRequest) (domain.Category, error) {
	name := strings.TrimSpace(req.Name)
	if id <= 0 || name == "" {
		return domain.Category{}, domain.ErrInvalidInput
	}
	return s.repo.UpdateCategory(ctx, id, name)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}

	count, err := s.repo.CategoryBookCount(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return domain.ErrCategoryInUse
	}

	return s.repo.DeleteCategory(ctx, id)
}
