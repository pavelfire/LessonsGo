package memory_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"lessonsgo/testmicroservice/internal/domain/port"
	"lessonsgo/testmicroservice/internal/storage/memory"
)

func TestStore_UpsertAndGet(t *testing.T) {
	ctx := context.Background()
	store := memory.New()

	first := port.Port{ID: "AEDXB", Name: "Dubai", Coordinates: []float64{55.27, 25.25}}
	second := port.Port{ID: "AEDXB", Name: "Dubai (updated)", Coordinates: []float64{55.2708, 25.2048}}

	require.NoError(t, store.Upsert(ctx, first))
	require.NoError(t, store.Upsert(ctx, second))

	got, ok := store.Get(ctx, "AEDXB")
	require.True(t, ok)
	assert.Equal(t, second, got)
	assert.Equal(t, 1, store.Count())
}

func TestStore_UpsertRespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	store := memory.New()
	err := store.Upsert(ctx, port.Port{ID: "NLRTM", Name: "Rotterdam"})

	require.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, 0, store.Count())
}
