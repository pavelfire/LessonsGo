package port

import "context"

type Repository interface {
	Upsert(ctx context.Context, p Port) error
	Get(ctx context.Context, id string) (Port, bool)
	Count() int
}
