package banksim

import (
	"context"
	"testing"
	"time"

	"lessonsgo/testpaymentgeteway/internal/domain/payment"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimulatorDeclinesHighAmount(t *testing.T) {
	sim := New()
	sim.now = func() time.Time {
		return time.Date(2026, 6, 13, 12, 0, 0, 0, time.UTC)
	}

	resp, err := sim.ProcessPayment(context.Background(), payment.BankRequest{
		CardNumber:  "4111111111111111",
		ExpiryMonth: 12,
		ExpiryYear:  2027,
		Amount:      100_001,
		Currency:    "USD",
	})
	require.NoError(t, err)
	assert.Equal(t, payment.StatusFailed, resp.Status)
	assert.NotEmpty(t, resp.ID)
}

func TestSimulatorDeclinesCardEnding0000(t *testing.T) {
	sim := New()
	sim.now = func() time.Time {
		return time.Date(2026, 6, 13, 12, 0, 0, 0, time.UTC)
	}

	resp, err := sim.ProcessPayment(context.Background(), payment.BankRequest{
		CardNumber:  "4111111111110000",
		ExpiryMonth: 12,
		ExpiryYear:  2027,
		Amount:      1000,
		Currency:    "USD",
	})
	require.NoError(t, err)
	assert.Equal(t, payment.StatusFailed, resp.Status)
}

func TestSimulatorApprovesValidPayment(t *testing.T) {
	sim := New()
	sim.now = func() time.Time {
		return time.Date(2026, 6, 13, 12, 0, 0, 0, time.UTC)
	}

	resp, err := sim.ProcessPayment(context.Background(), payment.BankRequest{
		CardNumber:  "4111111111111111",
		ExpiryMonth: 12,
		ExpiryYear:  2027,
		Amount:      1000,
		Currency:    "USD",
	})
	require.NoError(t, err)
	assert.Equal(t, payment.StatusSucceeded, resp.Status)
}
