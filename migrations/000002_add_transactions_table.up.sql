CREATE TABLE IF NOT EXISTS transactions
(
    transaction_id BIGSERIAL PRIMARY KEY,
    user_id        BIGINT         NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,

    operation_date TIMESTAMPTZ    NOT NULL,
    account_name   TEXT           NOT NULL,
    card_name      TEXT,
    amount         DECIMAL(12, 2) NOT NULL,
    category       TEXT           NOT NULL,
    operation_type TEXT           NOT NULL,
    comment        TEXT,
    bonus_value    DECIMAL(12, 2) DEFAULT 0,

    created_at     TIMESTAMPTZ    DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions (user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions (operation_date);