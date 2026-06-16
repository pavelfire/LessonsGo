package cart

import (
	"context"
	"time"

	"lessonsgo/testbookapi/internal/domain"
)

const ReservationTTL = 30 * time.Minute

type Service struct {
	repo domain.Repository
	now  func() time.Time
}

func NewService(repo domain.Repository) *Service {
	return &Service{
		repo: repo,
		now:  time.Now,
	}
}

type ItemResponse struct {
	BookID    int64     `json:"book_id"`
	AddedAt   time.Time `json:"added_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type CartResponse struct {
	Items []ItemResponse `json:"items"`
}

func (s *Service) List(ctx context.Context, userID int64) (CartResponse, error) {
	if err := s.releaseExpired(ctx); err != nil {
		return CartResponse{}, err
	}

	items, err := s.repo.ListCartItems(ctx, userID)
	if err != nil {
		return CartResponse{}, err
	}

	resp := CartResponse{Items: make([]ItemResponse, 0, len(items))}
	for _, item := range items {
		resp.Items = append(resp.Items, ItemResponse{
			BookID:    item.BookID,
			AddedAt:   item.AddedAt,
			ExpiresAt: item.ExpiresAt,
		})
	}
	return resp, nil
}

func (s *Service) Add(ctx context.Context, userID, bookID int64) error {
	if userID <= 0 || bookID <= 0 {
		return domain.ErrInvalidInput
	}

	if err := s.releaseExpired(ctx); err != nil {
		return err
	}

	book, err := s.repo.GetBook(ctx, bookID)
	if err != nil {
		return err
	}
	if book.Stock <= 0 {
		return domain.ErrOutOfStock
	}

	expiresAt := s.now().Add(ReservationTTL)
	return s.repo.AddCartItem(ctx, userID, bookID, expiresAt)
}

func (s *Service) Remove(ctx context.Context, userID, bookID int64) error {
	if userID <= 0 || bookID <= 0 {
		return domain.ErrInvalidInput
	}
	return s.repo.RemoveCartItem(ctx, userID, bookID)
}

func (s *Service) Checkout(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return domain.ErrInvalidInput
	}

	if err := s.releaseExpired(ctx); err != nil {
		return err
	}

	items, err := s.repo.ListCartItems(ctx, userID)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return domain.ErrCartEmpty
	}

	return s.repo.CheckoutCart(ctx, userID)
}

func (s *Service) ReleaseExpired(ctx context.Context) (int, error) {
	return s.repo.ReleaseExpiredCartItems(ctx)
}

func (s *Service) releaseExpired(ctx context.Context) error {
	_, err := s.repo.ReleaseExpiredCartItems(ctx)
	return err
}
