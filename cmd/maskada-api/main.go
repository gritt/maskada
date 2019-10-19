package main

import (
	"log"
	"net/http"
)

func main() {
	api, err := initAPI()
	if err != nil {
		log.Fatalln(err)
	}

	server := &http.Server{
		Addr:    ":8888",
		Handler: api.Routes(),
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}
