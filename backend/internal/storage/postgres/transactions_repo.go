package postgres

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"
)

func calculateHash(userID int64, bonus float64, opType string, date time.Time, amount float64, comment string) string {
	input := fmt.Sprintf("%d-%s-%f-%s-%f-%s", userID, date, amount, comment, bonus, opType)

	hash := sha256.Sum256([]byte(input))

	return fmt.Sprintf("%x", hash)
}

type Transaction struct {
	UserID        int64
	OperationDate time.Time
	AccountName   string
	CardName      string
	Amount        float64
	Category      string
	OperationType string
	Comment       string
	Bonus         float64
	CreatedAt     time.Time
}

func (db *DB) SaveTransaction(ctx context.Context, tx Transaction) error {

	txHash := calculateHash(tx.UserID, tx.Bonus, tx.OperationType, tx.OperationDate, tx.Amount, tx.Comment)

	query := `
        INSERT INTO transactions (
            user_id, 
            operation_date, 
            account_name, 
            card_name, 
            amount, 
            category, 
            operation_type, 
            comment, 
            bonus_value,
            transaction_hash
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        ON CONFLICT (transaction_hash) DO NOTHING;
    `

	_, err := db.Pool.Exec(ctx, query,
		tx.UserID,
		tx.OperationDate,
		tx.AccountName,
		tx.CardName,
		tx.Amount,
		tx.Category,
		tx.OperationType,
		tx.Comment,
		tx.Bonus,
		txHash,
	)

	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	return nil
}
