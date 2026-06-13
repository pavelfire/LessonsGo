package payment

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessRequestValidate(t *testing.T) {
	now := time.Date(2026, 6, 13, 12, 0, 0, 0, time.UTC)

	valid := ProcessRequest{
		CardNumber:  "4111111111111111",
		ExpiryMonth: 12,
		ExpiryYear:  2027,
		CVV:         "123",
		Amount:      1000,
		Currency:    "usd",
	}
	require.NoError(t, valid.Validate(now))

	tests := []struct {
		name string
		req  ProcessRequest
		err  error
	}{
		{
			name: "invalid luhn",
			req: ProcessRequest{
				CardNumber:  "4111111111111112",
				ExpiryMonth: 12,
				ExpiryYear:  2027,
				CVV:         "123",
				Amount:      1000,
				Currency:    "USD",
			},
			err: ErrInvalidCardNumber,
		},
		{
			name: "expired card",
			req: ProcessRequest{
				CardNumber:  "4111111111111111",
				ExpiryMonth: 1,
				ExpiryYear:  2020,
				CVV:         "123",
				Amount:      1000,
				Currency:    "USD",
			},
			err: ErrInvalidExpiry,
		},
		{
			name: "invalid cvv",
			req: ProcessRequest{
				CardNumber:  "4111111111111111",
				ExpiryMonth: 12,
				ExpiryYear:  2027,
				CVV:         "12",
				Amount:      1000,
				Currency:    "USD",
			},
			err: ErrInvalidCVV,
		},
		{
			name: "zero amount",
			req: ProcessRequest{
				CardNumber:  "4111111111111111",
				ExpiryMonth: 12,
				ExpiryYear:  2027,
				CVV:         "123",
				Amount:      0,
				Currency:    "USD",
			},
			err: ErrInvalidAmount,
		},
		{
			name: "invalid currency",
			req: ProcessRequest{
				CardNumber:  "4111111111111111",
				ExpiryMonth: 12,
				ExpiryYear:  2027,
				CVV:         "123",
				Amount:      1000,
				Currency:    "RUB",
			},
			err: ErrInvalidCurrency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate(now)
			require.Error(t, err)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func TestMaskCardNumber(t *testing.T) {
	assert.Equal(t, "**** **** **** 1111", MaskCardNumber("4111-1111-1111-1111"))
	assert.Equal(t, "1111", CardLast4("4111111111111111"))
}

func TestValidLuhn(t *testing.T) {
	assert.True(t, validLuhn("4111111111111111"))
	assert.False(t, validLuhn("4111111111111112"))
}
