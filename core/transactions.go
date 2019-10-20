package core

import "github.com/pkg/errors"

type (
	// Repository represents a client able to save a transaction.
	Repository interface {
		Create(Transaction) (Transaction, error)
	}

	// CreateTransactionUseCase implements the business logic to create a transaction.
	CreateTransactionUseCase struct {
		repository Repository
	}
)

// NewCreateTransactionUseCase initialize the use case.
func NewCreateTransactionUseCase(r Repository) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{repository: r}
}

// Create a transaction.
func (uc *CreateTransactionUseCase) Create(t Transaction) (Transaction, error) {
	if err := t.Validate(); err != nil {
		return Transaction{}, errors.Wrap(err, "Create failed")
	}

	transaction, err := uc.repository.Create(t)
	if err != nil {
		return Transaction{}, errors.Wrap(err, "Create failed")
	}

	return transaction, nil
}
