package book

import (
	"context"
	"strings"

	"lessonsgo/testbookapi/internal/domain"
)

const (
	defaultLimit = 50
	maxLimit     = 100
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

type CreateRequest struct {
	Title      string  `json:"title"`
	Year       int     `json:"year_published"`
	Author     string  `json:"author"`
	PriceUSD   float64 `json:"price_usd"`
	CategoryID int64   `json:"category_id"`
	Stock      int     `json:"stock"`
}

type UpdateRequest struct {
	Title      string  `json:"title"`
	Year       int     `json:"year_published"`
	Author     string  `json:"author"`
	PriceUSD   float64 `json:"price_usd"`
	CategoryID int64   `json:"category_id"`
	Stock      *int    `json:"stock,omitempty"`
}

type ListResult struct {
	Books  []domain.Book `json:"books"`
	Total  int           `json:"total"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (domain.Book, error) {
	if err := validateBookFields(req.Title, req.Year, req.Author, req.PriceUSD, req.CategoryID); err != nil {
		return domain.Book{}, err
	}
	if req.Stock < 0 {
		return domain.Book{}, domain.ErrInvalidInput
	}

	if _, err := s.repo.GetCategory(ctx, req.CategoryID); err != nil {
		return domain.Book{}, err
	}

	return s.repo.CreateBook(ctx, domain.Book{
		Title:      strings.TrimSpace(req.Title),
		Year:       req.Year,
		Author:     strings.TrimSpace(req.Author),
		PriceUSD:   req.PriceUSD,
		CategoryID: req.CategoryID,
		Stock:      req.Stock,
	})
}

func (s *Service) Get(ctx context.Context, id int64) (domain.Book, error) {
	if id <= 0 {
		return domain.Book{}, domain.ErrInvalidInput
	}

	book, err := s.repo.GetBook(ctx, id)
	if err != nil {
		return domain.Book{}, err
	}
	if book.Stock <= 0 {
		return domain.Book{}, domain.ErrNotFound
	}
	return book, nil
}

func (s *Service) List(ctx context.Context, categoryIDs []int64, limit, offset int) (ListResult, error) {
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}

	filter := domain.BookFilter{
		CategoryIDs: categoryIDs,
		Limit:       limit,
		Offset:      offset,
	}

	books, err := s.repo.ListBooks(ctx, filter)
	if err != nil {
		return ListResult{}, err
	}

	total, err := s.repo.CountBooks(ctx, filter)
	if err != nil {
		return ListResult{}, err
	}

	return ListResult{
		Books:  books,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *Service) Update(ctx context.Context, id int64, req UpdateRequest) (domain.Book, error) {
	if req.Stock != nil {
		return domain.Book{}, domain.ErrStockNotEditable
	}
	if err := validateBookFields(req.Title, req.Year, req.Author, req.PriceUSD, req.CategoryID); err != nil {
		return domain.Book{}, err
	}
	if id <= 0 {
		return domain.Book{}, domain.ErrInvalidInput
	}

	if _, err := s.repo.GetCategory(ctx, req.CategoryID); err != nil {
		return domain.Book{}, err
	}

	return s.repo.UpdateBook(ctx, id,
		strings.TrimSpace(req.Title),
		req.Year,
		strings.TrimSpace(req.Author),
		req.PriceUSD,
		req.CategoryID,
	)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}
	return s.repo.DeleteBook(ctx, id)
}

func validateBookFields(title string, year int, author string, price float64, categoryID int64) error {
	if strings.TrimSpace(title) == "" ||
		strings.TrimSpace(author) == "" ||
		year < 0 ||
		price < 0 ||
		categoryID <= 0 {
		return domain.ErrInvalidInput
	}
	return nil
}
