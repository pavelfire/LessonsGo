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
