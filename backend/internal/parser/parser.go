package parser

import (
	"backend/internal/storage/postgres"
	"backend/pkg/encoding"
	"strconv"
	"strings"
	"time"
)

type StatementParser interface {
	Parse(record []string) (*postgres.Transaction, error)
	GetDelimiter() rune
}

type AlfaParser struct{}

type TBankParser struct{}

func (p *AlfaParser) GetDelimiter() rune { return ';' }

func (p *AlfaParser) Parse(record []string) (*postgres.Transaction, error) {
	amount, _ := strconv.ParseFloat(record[7], 64)
	opDate, _ := time.Parse("02.01.2006", record[0])
	bonus, _ := strconv.ParseFloat(record[14], 64)

	rawAccountName := record[2]
	accountName := encoding.SafeDecode(rawAccountName)

	rawCategory := record[10]
	category := encoding.SafeDecode(rawCategory)

	rawOperation := record[12]
	operation := encoding.SafeDecode(rawOperation)

	rawComment := record[13]
	comment := encoding.SafeDecode(rawComment)

	rawCardName := record[4]
	cardName := encoding.SafeDecode(rawCardName)

	return &postgres.Transaction{
		OperationDate: opDate,
		AccountName:   accountName,
		CardName:      cardName,
		Amount:        amount,
		Category:      category,
		OperationType: operation,
		Comment:       comment,
		Bonus:         bonus,
	}, nil
}

func (t *TBankParser) GetDelimiter() rune { return ';' }
func (t *TBankParser) Parse(record []string) (*postgres.Transaction, error) {
	rawAmount := record[4]
	rawAmount = strings.Trim(rawAmount, "\"")
	rawAmount = strings.ReplaceAll(rawAmount, ",", ".")
	amount, _ := strconv.ParseFloat(rawAmount, 64)
	opDate, _ := time.Parse("02.01.2006 15:04", record[0])
	bonus, _ := strconv.ParseFloat(record[12], 64)

	var opType string
	if amount < 0 {
		opType = "Списание"
	} else {
		opType = "Пополнение"
	}

	rawCategory := record[9]
	category := encoding.SafeDecode(rawCategory)

	return &postgres.Transaction{
		OperationDate: opDate,
		AccountName:   "TBank",
		CardName:      "TBank",
		Amount:        amount,
		Category:      category,
		OperationType: opType,
		Comment:       "",
		Bonus:         bonus,
	}, nil
}
