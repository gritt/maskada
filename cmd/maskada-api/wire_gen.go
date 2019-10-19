// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/gritt/maskada/core"
	"github.com/gritt/maskada/details"
	"github.com/gritt/maskada/details/db"
	"github.com/gritt/maskada/details/rest"
)

// Injectors from wire.go:

func initAPI() (*rest.API, error) {
	config, err := details.NewConfig()
	if err != nil {
		return nil, err
	}
	repository, err := db.NewRepository(config)
	if err != nil {
		return nil, err
	}
	createTransactionUseCase := core.NewCreateTransactionUseCase(repository)
	api := rest.NewAPI(createTransactionUseCase)
	return api, nil
}

// wire.go:

var repositorySet = wire.NewSet(details.NewConfig, wire.Bind(new(core.Repository), new(*db.Repository)), db.NewRepository)

var createTransactionSet = wire.NewSet(wire.Bind(new(rest.TransactionCreator), new(*core.CreateTransactionUseCase)), core.NewCreateTransactionUseCase)
