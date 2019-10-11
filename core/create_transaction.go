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

// TODO @gritt: Test this!
func NewCreateTransactionUseCase(r Repository) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		repository: r,
	}
}

func (uc *CreateTransactionUseCase) Create(t Transaction) (Transaction, error) {

	if err := t.Validate(); err != nil {
		return Transaction{}, errors.Wrap(err, "Create failed")
	}

	transaction, err := uc.repository.Create(t)
	if err != nil {
		return Transaction{}, errors.Wrap(err, "Create failed")
	}

	return transaction, err
}
