package banksimapi

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"lessonsgo/testpaymentgeteway/internal/banksim"
	"lessonsgo/testpaymentgeteway/internal/domain/payment"
)

type Handler struct {
	sim    *banksim.Simulator
	logger *slog.Logger
}

func NewHandler(sim *banksim.Simulator, logger *slog.Logger) *Handler {
	return &Handler{sim: sim, logger: logger}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("POST /process", h.process)
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
	ID     string         `json:"id"`
	Status payment.Status `json:"status"`
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (h *Handler) process(w http.ResponseWriter, r *http.Request) {
	var req processRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	resp, err := h.sim.ProcessPayment(r.Context(), payment.BankRequest{
		CardNumber:  req.CardNumber,
		ExpiryMonth: req.ExpiryMonth,
		ExpiryYear:  req.ExpiryYear,
		CVV:         req.CVV,
		Amount:      req.Amount,
		Currency:    req.Currency,
	})
	if err != nil {
		h.logger.Error("bank processing failed", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, processResponse{
		ID:     resp.ID,
		Status: resp.Status,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
