package bank

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"lessonsgo/testpaymentgeteway/internal/domain/payment"
)

type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: stringsTrimRightSlash(baseURL),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type processRequest struct {
	CardNumber  string `json:"card_number"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	CVV         string `json:"cvv"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
}

type processResponse struct {
	ID     string          `json:"id"`
	Status payment.Status  `json:"status"`
}

func (c *HTTPClient) ProcessPayment(ctx context.Context, req payment.BankRequest) (payment.BankResponse, error) {
	body, err := json.Marshal(processRequest{
		CardNumber:  req.CardNumber,
		ExpiryMonth: req.ExpiryMonth,
		ExpiryYear:  req.ExpiryYear,
		CVV:         req.CVV,
		Amount:      req.Amount,
		Currency:    req.Currency,
	})
	if err != nil {
		return payment.BankResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/process", bytes.NewReader(body))
	if err != nil {
		return payment.BankResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return payment.BankResponse{}, fmt.Errorf("bank request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return payment.BankResponse{}, fmt.Errorf("bank returned status %d: %s", resp.StatusCode, string(msg))
	}

	var out processResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return payment.BankResponse{}, fmt.Errorf("decode bank response: %w", err)
	}

	return payment.BankResponse{
		ID:     out.ID,
		Status: out.Status,
	}, nil
}

func stringsTrimRightSlash(s string) string {
	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}
