package core

import "github.com/pkg/errors"

type (
	// Repository represents a client able to save and find a transaction.
	Repository interface {
		Create(Transaction) (Transaction, error)
		Find() ([]Transaction, error)
	}

	// CreateTransactionUseCase implements the business logic to create a transaction.
	CreateTransactionUseCase struct {
		repository Repository
	}

	// ListTransactionUseCase implements the business logic to find a transaction.
	ListTransactionUseCase struct {
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

// NewListTransactionUseCase initialize the use case.
func NewListTransactionUseCase(r Repository) *ListTransactionUseCase {
	return &ListTransactionUseCase{repository: r}
}

// List transaction(s).
func (uc *ListTransactionUseCase) List() ([]Transaction, error) {
	transactions, err := uc.repository.Find()
	if err != nil {
		return []Transaction{}, errors.Wrap(err, "List failed")
	}

	return transactions, nil
}
