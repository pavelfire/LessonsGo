package port_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"lessonsgo/testmicroservice/internal/domain/port"
)

type mockRepository struct {
	upsert func(context.Context, port.Port) error
	get    func(context.Context, string) (port.Port, bool)
	count  int
}

func (m *mockRepository) Upsert(ctx context.Context, p port.Port) error {
	if m.upsert != nil {
		return m.upsert(ctx, p)
	}
	return nil
}

func (m *mockRepository) Get(ctx context.Context, id string) (port.Port, bool) {
	if m.get != nil {
		return m.get(ctx, id)
	}
	return port.Port{}, false
}

func (m *mockRepository) Count() int {
	return m.count
}

func TestService_Upsert(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("rejects port without id", func(t *testing.T) {
		t.Parallel()

		svc := port.NewService(&mockRepository{})
		err := svc.Upsert(ctx, port.Port{Name: "Dubai"})

		require.ErrorIs(t, err, port.ErrInvalidID)
	})

	t.Run("delegates valid port to repository", func(t *testing.T) {
		t.Parallel()

		var stored port.Port
		repo := &mockRepository{
			upsert: func(_ context.Context, p port.Port) error {
				stored = p
				return nil
			},
		}
		svc := port.NewService(repo)

		input := port.Port{ID: "AEDXB", Name: "Dubai"}
		err := svc.Upsert(ctx, input)

		require.NoError(t, err)
		assert.Equal(t, input, stored)
	})
}
