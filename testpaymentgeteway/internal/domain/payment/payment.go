package payment

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidCardNumber = errors.New("invalid card number")
	ErrInvalidExpiry     = errors.New("invalid card expiry")
	ErrInvalidCVV        = errors.New("invalid cvv")
	ErrInvalidAmount     = errors.New("amount must be greater than zero")
	ErrInvalidCurrency   = errors.New("invalid currency code")
	ErrNotFound          = errors.New("payment not found")
)

type Status string

const (
	StatusSucceeded Status = "succeeded"
	StatusFailed    Status = "failed"
)

type Payment struct {
	ID          string    `json:"id"`
	Amount      int64     `json:"amount"`
	Currency    string    `json:"currency"`
	CardMasked  string    `json:"card_masked"`
	CardLast4   string    `json:"card_last4"`
	ExpiryMonth int       `json:"expiry_month"`
	ExpiryYear  int       `json:"expiry_year"`
	Status      Status    `json:"status"`
	BankRefID   string    `json:"bank_reference,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProcessRequest struct {
	CardNumber  string `json:"card_number"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	CVV         string `json:"cvv"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
}

type BankRequest struct {
	CardNumber  string
	ExpiryMonth int
	ExpiryYear  int
	CVV         string
	Amount      int64
	Currency    string
}

type BankResponse struct {
	ID     string
	Status Status
}

func (r ProcessRequest) Validate(now time.Time) error {
	card := normalizeCardNumber(r.CardNumber)
	if !validLuhn(card) {
		return ErrInvalidCardNumber
	}
	if r.ExpiryMonth < 1 || r.ExpiryMonth > 12 {
		return ErrInvalidExpiry
	}
	expiry := time.Date(r.ExpiryYear, time.Month(r.ExpiryMonth+1), 0, 23, 59, 59, 0, time.UTC)
	if expiry.Before(now) {
		return ErrInvalidExpiry
	}
	if len(r.CVV) < 3 || len(r.CVV) > 4 || !isDigits(r.CVV) {
		return ErrInvalidCVV
	}
	if r.Amount <= 0 {
		return ErrInvalidAmount
	}
	if !validCurrency(r.Currency) {
		return ErrInvalidCurrency
	}
	return nil
}

func (r ProcessRequest) ToBankRequest() BankRequest {
	return BankRequest{
		CardNumber:  normalizeCardNumber(r.CardNumber),
		ExpiryMonth: r.ExpiryMonth,
		ExpiryYear:  r.ExpiryYear,
		CVV:         r.CVV,
		Amount:      r.Amount,
		Currency:    strings.ToUpper(r.Currency),
	}
}

func MaskCardNumber(card string) string {
	digits := normalizeCardNumber(card)
	if len(digits) < 4 {
		return "****"
	}
	last4 := digits[len(digits)-4:]
	return fmt.Sprintf("**** **** **** %s", last4)
}

func CardLast4(card string) string {
	digits := normalizeCardNumber(card)
	if len(digits) < 4 {
		return digits
	}
	return digits[len(digits)-4:]
}

func normalizeCardNumber(card string) string {
	var b strings.Builder
	for _, r := range card {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}

func validCurrency(currency string) bool {
	switch strings.ToUpper(currency) {
	case "USD", "EUR", "GBP":
		return true
	default:
		return false
	}
}

func validLuhn(card string) bool {
	if len(card) < 13 || len(card) > 19 {
		return false
	}
	sum := 0
	alt := false
	for i := len(card) - 1; i >= 0; i-- {
		n := int(card[i] - '0')
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}
