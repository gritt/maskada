package core

import (
	"errors"
	"time"
)

const (
	// Debit is a transaction which is subtracted.
	Debit = 1

	// Credit is a transaction which is subtracted the next month.
	Credit = 2

	// Income is a transaction which is summed.
	Income = 3
)

type (
	// Category is the general class of a Transaction (eg: Health, Food).
	Category struct {
		Name string
	}

	// Transaction is money received or expended.
	Transaction struct {
		ID       int
		Amount   int
		Type     int
		Category Category
		Date     time.Time
		Name     string
	}
)

// Validate whether a transaction has all it's required properties set.
func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return errors.New("Transaction.Validate: missing amount")
	}

	if t.Type != Debit && t.Type != Credit && t.Type != Income {
		return errors.New("Transaction.Validate: missing type")
	}

	if t.Category.Name == "" {
		return errors.New("Transaction.Validate: missing category")
	}

	return nil
}
