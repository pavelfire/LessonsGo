package payment

import "context"

type Repository interface {
	Save(ctx context.Context, p Payment) error
	GetByID(ctx context.Context, id string) (Payment, bool)
}
