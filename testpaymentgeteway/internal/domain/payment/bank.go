package payment

import "context"

type BankClient interface {
	ProcessPayment(ctx context.Context, req BankRequest) (BankResponse, error)
}
