package rest

import (
	"net/http"

	"github.com/go-chi/chi"
)

// TODO @gritt: Test this!
func (api *API) Routes() *chi.Mux {
	r := chi.NewRouter()

	// TODO @gritt: Create middleware to handle authenticated users

	r.Route("/v1", func(r chi.Router) {
		r.Method(http.MethodPost, "/transaction", api.HandleCreateTransaction())
	})

	return r
}
