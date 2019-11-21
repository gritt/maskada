// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/gritt/maskada/core"
	"github.com/gritt/maskada/details"
	"github.com/gritt/maskada/details/db"
	"github.com/gritt/maskada/details/rest"
)

var repositorySet = wire.NewSet(
	details.NewConfig,
	wire.Bind(new(core.Repository), new(*db.Repository)),
	db.NewRepository,
)

var createTransactionSet = wire.NewSet(
	wire.Bind(new(rest.TransactionCreator), new(*core.CreateTransactionUseCase)),
	core.NewCreateTransactionUseCase,
)

var listTransactionSet = wire.NewSet(
	wire.Bind(new(rest.TransactionLister), new(*core.ListTransactionUseCase)),
	core.NewListTransactionUseCase,
)

func initAPI() (*rest.API, error) {
	panic(wire.Build(
		repositorySet,
		createTransactionSet,
		listTransactionSet,
		rest.NewAPI,
	))
}
