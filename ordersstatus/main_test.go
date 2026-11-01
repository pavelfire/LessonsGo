package domain

import (
	"context"
	"testing"
	"time"
)

func TestSaveOrderStatusHistoryValidation(t *testing.T) {
	store := NewStore(nil)
	now := time.Date(2026, 10, 30, 7, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		orderID   string
		status    string
		partition string
		at        time.Time
	}{
		{name: "empty order", status: "paid", partition: "2026-10", at: now},
		{name: "empty status", orderID: "order-1", partition: "2026-10", at: now},
		{name: "empty partition", orderID: "order-1", status: "paid", at: now},
		{name: "zero time", orderID: "order-1", status: "paid", partition: "2026-10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.SaveOrderStatusHistory(context.Background(), tt.orderID, tt.status, tt.partition, tt.at)
			if err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
