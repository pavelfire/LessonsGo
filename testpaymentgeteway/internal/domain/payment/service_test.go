package payment

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRepo struct {
	saved Payment
}

func (m *mockRepo) Save(_ context.Context, p Payment) error {
	m.saved = p
	return nil
}

func (m *mockRepo) GetByID(_ context.Context, id string) (Payment, bool) {
	if m.saved.ID == id {
		return m.saved, true
	}
	return Payment{}, false
}

type mockBank struct {
	resp BankResponse
	err  error
}

func (m *mockBank) ProcessPayment(_ context.Context, _ BankRequest) (BankResponse, error) {
	return m.resp, m.err
}

func TestServiceProcess(t *testing.T) {
	repo := &mockRepo{}
	bank := &mockBank{
		resp: BankResponse{ID: "bank_abc", Status: StatusSucceeded},
	}

	svc := NewService(repo, bank)
	svc.now = func() time.Time {
		return time.Date(2026, 6, 13, 12, 0, 0, 0, time.UTC)
	}

	p, err := svc.Process(context.Background(), ProcessRequest{
		CardNumber:  "4111111111111111",
		ExpiryMonth: 12,
		ExpiryYear:  2027,
		CVV:         "123",
		Amount:      2500,
		Currency:    "USD",
	})
	require.NoError(t, err)
	assert.Equal(t, StatusSucceeded, p.Status)
	assert.Equal(t, "bank_abc", p.BankRefID)
	assert.Equal(t, "**** **** **** 1111", p.CardMasked)
	assert.Equal(t, int64(2500), repo.saved.Amount)
}

func TestServiceGetNotFound(t *testing.T) {
	svc := NewService(&mockRepo{}, &mockBank{})
	_, err := svc.Get(context.Background(), "missing")
	assert.ErrorIs(t, err, ErrNotFound)
}
