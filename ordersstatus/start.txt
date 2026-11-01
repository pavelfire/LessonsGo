/*
Есть заказы, которые в процессе жизненного цикла переходят по статусам.
Требуется реализовать слой работы с БД:
- Сохранять историю статусов заказа
- Учесть партицирование /
В дальнейшем, за рамками текущей задачи, будут реализованы другие методы: получение статусов по времени, полная история заказа и т.п.
*/
package domain

import (
"database/sql"
"fmt"
)

type OrderStatusHistoryStore struct {}

func NewStore() OrderStatusHistoryStore {
	return OrderStatusHistoryStore{}
}

func (s OrderStatusHistoryStore) saveOrderStatusHistory(orderId, status, partition string) {
	db, err := sql. Open("postgres", "user=foo dbname=bar sslmode=disable")
	if err != nil {
		fmt. Errorf("%v", err)
	}
	Query := fmt.Sprintf("insert into orderstatushistory (order_id, status, partition, inserted_at) values (%v, %v, %v, %v)", orderId, partition, status, time. Now())
	go db. Exec(Query)
}