package rest

import (
	"github.com/gritt/maskada/core"
)

// TransactionCreator represents a use case able to create a transaction.
type TransactionCreator interface {
	Create(core.Transaction) (core.Transaction, error)
}

// API holds all use cases.
type API struct {
	TransactionCreator TransactionCreator
}

// NewAPI initialize the API.
func NewAPI(creator TransactionCreator) *API {
	return &API{
		TransactionCreator: creator,
	}
}
