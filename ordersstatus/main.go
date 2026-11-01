/*
Есть заказы, которые в процессе жизненного цикла переходят по статусам.
Требуется реализовать слой работы с БД:
- Сохранять историю статусов заказа
- Учесть партицирование /
В дальнейшем, за рамками текущей задачи, будут реализованы другие методы: получение статусов по времени, полная история заказа и т.п.

Таблица order_status_history партиционирована по LIST (partition).
INSERT идёт в родительскую таблицу, PostgreSQL направляет строку в партицию.
Партиции создаются заранее, ключ партиции входит в первичный ключ.
*/
package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type OrderStatusHistoryStore struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *OrderStatusHistoryStore {
	return &OrderStatusHistoryStore{db: db}
}

func (s *OrderStatusHistoryStore) SaveOrderStatusHistory(ctx context.Context, orderID, status, partition string, insertedAt time.Time) error {
	if orderID == "" || status == "" || partition == "" {
		return errors.New("orderID, status and partition are required")
	}
	if insertedAt.IsZero() {
		return errors.New("insertedAt is required")
	}
	if s == nil || s.db == nil {
		return errors.New("order status history store is not initialized")
	}

	const query = `
		INSERT INTO order_status_history (order_id, status, partition, inserted_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := s.db.ExecContext(ctx, query, orderID, status, partition, insertedAt.UTC())
	if err != nil {
		return fmt.Errorf("save order status history: %w", err)
	}
	return nil
}


/*
CREATE TABLE order_status_history (
    order_id    text        NOT NULL,
    status      text        NOT NULL,
    partition   text        NOT NULL,
    inserted_at timestamptz NOT NULL,
    PRIMARY KEY (partition, inserted_at, order_id)
) PARTITION BY LIST (partition);

CREATE TABLE order_status_history_2026_10
    PARTITION OF order_status_history
    FOR VALUES IN ('2026-10');
*/