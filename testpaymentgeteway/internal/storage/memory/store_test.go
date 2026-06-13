package memory

import (
	"context"
	"testing"

	"lessonsgo/testpaymentgeteway/internal/domain/payment"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreSaveAndGet(t *testing.T) {
	store := New()
	p := payment.Payment{
		ID:         "pay_1",
		Amount:     1000,
		Currency:   "USD",
		CardMasked: "**** **** **** 1111",
		Status:     payment.StatusSucceeded,
	}

	require.NoError(t, store.Save(context.Background(), p))

	got, ok := store.GetByID(context.Background(), "pay_1")
	require.True(t, ok)
	assert.Equal(t, p, got)
	assert.Equal(t, 1, store.Count())
}
