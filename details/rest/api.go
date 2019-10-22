package rest

import (
	"github.com/gritt/maskada/core"
)

type (
	// TransactionCreator represents a use case able to create a transaction.
	TransactionCreator interface {
		Create(core.Transaction) (core.Transaction, error)
	}

	// TransactionLister represents a use case able to list transactions.
	TransactionLister interface {
		List() ([]core.Transaction, error)
	}
)

// API holds all use cases.
type API struct {
	TransactionCreator TransactionCreator
	TransactionLister  TransactionLister
}

// NewAPI initialize the API.
func NewAPI(creator TransactionCreator, lister TransactionLister) *API {
	return &API{
		TransactionCreator: creator,
		TransactionLister:  lister,
	}
}
