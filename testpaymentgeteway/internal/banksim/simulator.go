package banksim

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"lessonsgo/testpaymentgeteway/internal/domain/payment"
)

const maxAmountMinor = 100_000 // $1000.00 in cents

type Simulator struct {
	now func() time.Time
}

func New() *Simulator {
	return &Simulator{now: time.Now}
}

func (s *Simulator) ProcessPayment(_ context.Context, req payment.BankRequest) (payment.BankResponse, error) {
	status := payment.StatusSucceeded

	if req.Amount > maxAmountMinor {
		status = payment.StatusFailed
	}
	if strings.HasSuffix(req.CardNumber, "0000") {
		status = payment.StatusFailed
	}
	expiry := time.Date(req.ExpiryYear, time.Month(req.ExpiryMonth+1), 0, 23, 59, 59, 0, time.UTC)
	if expiry.Before(s.now()) {
		status = payment.StatusFailed
	}

	return payment.BankResponse{
		ID:     newBankRefID(),
		Status: status,
	}, nil
}

func newBankRefID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("bank_%s", hex.EncodeToString(b))
}
