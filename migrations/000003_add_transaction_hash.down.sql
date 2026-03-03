DROP INDEX IF EXISTS idx_transactions_hash;
ALTER TABLE transactions DROP COLUMN IF EXISTS transaction_hash;