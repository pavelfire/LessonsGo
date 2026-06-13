package payment

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type Service struct {
	repo Repository
	bank BankClient
	now  func() time.Time
}

func NewService(repo Repository, bank BankClient) *Service {
	return &Service{
		repo: repo,
		bank: bank,
		now:  time.Now,
	}
}

func (s *Service) Process(ctx context.Context, req ProcessRequest) (Payment, error) {
	if err := req.Validate(s.now()); err != nil {
		return Payment{}, err
	}

	bankResp, err := s.bank.ProcessPayment(ctx, req.ToBankRequest())
	if err != nil {
		return Payment{}, fmt.Errorf("bank processing: %w", err)
	}

	p := Payment{
		ID:          newPaymentID(),
		Amount:      req.Amount,
		Currency:    req.ToBankRequest().Currency,
		CardMasked:  MaskCardNumber(req.CardNumber),
		CardLast4:   CardLast4(req.CardNumber),
		ExpiryMonth: req.ExpiryMonth,
		ExpiryYear:  req.ExpiryYear,
		Status:      bankResp.Status,
		BankRefID:   bankResp.ID,
		CreatedAt:   s.now().UTC(),
	}

	if err := s.repo.Save(ctx, p); err != nil {
		return Payment{}, err
	}

	return p, nil
}

func (s *Service) Get(ctx context.Context, id string) (Payment, error) {
	p, ok := s.repo.GetByID(ctx, id)
	if !ok {
		return Payment{}, ErrNotFound
	}
	return p, nil
}

func newPaymentID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "pay_" + hex.EncodeToString(b)
}
