package core

import (
	"errors"
	"time"
)

const (
	Debit  = 1
	Credit = 2
	Income = 3
)

type (
	Category struct {
		Name string
	}

	Transaction struct {
		ID       int
		Amount   int
		Type     int
		Category Category
		Date     time.Time
		Name     string
	}
)

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
