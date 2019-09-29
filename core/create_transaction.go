package core

import "github.com/pkg/errors"

type (
	Repository interface {
		Create(Transaction) (Transaction, error)
	}

	CreateTransactionUseCase struct {
		repository Repository
	}
)

func (uc CreateTransactionUseCase) CreateTransaction(t Transaction) (Transaction, error) {

	if err := t.Validate(); err != nil {
		return Transaction{}, errors.Wrap(err, "CreateTransaction failed")
	}

	transaction, err := uc.repository.Create(t)
	if err != nil {
		return Transaction{}, errors.Wrap(err, "CreateTransaction failed")
	}

	return transaction, err
}
