ALTER TABLE transactions ADD COLUMN IF NOT EXISTS transaction_hash VARCHAR(64);

CREATE UNIQUE INDEX IF NOT EXISTS idx_transactions_hash ON transactions (transaction_hash);