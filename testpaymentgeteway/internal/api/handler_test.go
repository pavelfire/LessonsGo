package api_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"lessonsgo/testpaymentgeteway/internal/api"
	"lessonsgo/testpaymentgeteway/internal/banksim"
	"lessonsgo/testpaymentgeteway/internal/domain/payment"
	"lessonsgo/testpaymentgeteway/internal/storage/memory"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentAPI(t *testing.T) {
	store := memory.New()
	svc := payment.NewService(store, banksim.New())

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	handler := api.NewHandler(svc, logger)

	mux := http.NewServeMux()
	handler.Register(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	body := map[string]any{
		"card_number":  "4111111111111111",
		"expiry_month": 12,
		"expiry_year":  2027,
		"cvv":          "123",
		"amount":       1500,
		"currency":     "USD",
	}
	payload, err := json.Marshal(body)
	require.NoError(t, err)

	createResp, err := http.Post(server.URL+"/payments", "application/json", bytes.NewReader(payload))
	require.NoError(t, err)
	defer createResp.Body.Close()
	require.Equal(t, http.StatusCreated, createResp.StatusCode)

	var created map[string]any
	require.NoError(t, json.NewDecoder(createResp.Body).Decode(&created))
	id, ok := created["id"].(string)
	require.True(t, ok)
	assert.Equal(t, "succeeded", created["status"])
	assert.Equal(t, "**** **** **** 1111", created["card_masked"])

	getResp, err := http.Get(server.URL + "/payments/" + id)
	require.NoError(t, err)
	defer getResp.Body.Close()
	require.Equal(t, http.StatusOK, getResp.StatusCode)

	var got map[string]any
	require.NoError(t, json.NewDecoder(getResp.Body).Decode(&got))
	assert.Equal(t, id, got["id"])
	assert.Equal(t, "succeeded", got["status"])
}
