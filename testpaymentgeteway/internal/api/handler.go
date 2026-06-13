package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"lessonsgo/testpaymentgeteway/internal/domain/payment"
)

type Handler struct {
	svc    *payment.Service
	logger *slog.Logger
}

func NewHandler(svc *payment.Service, logger *slog.Logger) *Handler {
	return &Handler{svc: svc, logger: logger}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /payments", h.processPayment)
	mux.HandleFunc("GET /payments/{id}", h.getPayment)
}

type processPaymentResponse struct {
	ID          string          `json:"id"`
	Status      payment.Status  `json:"status"`
	BankRefID   string          `json:"bank_reference,omitempty"`
	Amount      int64           `json:"amount"`
	Currency    string          `json:"currency"`
	CardMasked  string          `json:"card_masked"`
	ExpiryMonth int             `json:"expiry_month"`
	ExpiryYear  int             `json:"expiry_year"`
	CreatedAt   string          `json:"created_at"`
}

func (h *Handler) processPayment(w http.ResponseWriter, r *http.Request) {
	var req payment.ProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	p, err := h.svc.Process(r.Context(), req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toProcessResponse(p))
}

func (h *Handler) getPayment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "payment id is required")
		return
	}

	p, err := h.svc.Get(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toProcessResponse(p))
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, payment.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, payment.ErrInvalidCardNumber),
		errors.Is(err, payment.ErrInvalidExpiry),
		errors.Is(err, payment.ErrInvalidCVV),
		errors.Is(err, payment.ErrInvalidAmount),
		errors.Is(err, payment.ErrInvalidCurrency):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		h.logger.Error("request failed", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func toProcessResponse(p payment.Payment) processPaymentResponse {
	return processPaymentResponse{
		ID:          p.ID,
		Status:      p.Status,
		BankRefID:   p.BankRefID,
		Amount:      p.Amount,
		Currency:    p.Currency,
		CardMasked:  p.CardMasked,
		ExpiryMonth: p.ExpiryMonth,
		ExpiryYear:  p.ExpiryYear,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
