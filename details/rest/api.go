package rest

import (
	"github.com/gritt/maskada/core"
)

type TransactionCreator interface {
	Create(core.Transaction) (core.Transaction, error)
}

type API struct {
	TransactionCreator TransactionCreator
}

func NewAPI(creator TransactionCreator) *API {
	return &API{
		TransactionCreator: creator,
	}
}
