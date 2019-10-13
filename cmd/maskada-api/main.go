package main

import (
	"log"
	"net/http"

	"github.com/gritt/maskada/core"
	"github.com/gritt/maskada/details"
	"github.com/gritt/maskada/details/db"
	"github.com/gritt/maskada/details/rest"
)

func main() {
	// TODO @gritt: Test this!
	cfg, err := details.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	r, err := db.NewRepository(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	uc := core.NewCreateTransactionUseCase(r)

	// TODO @gritt: Use Dependecy Injection
	api := rest.NewAPI(uc)

	server := &http.Server{
		Addr:    ":8888",
		Handler: api.Routes(),
	}

	if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}
