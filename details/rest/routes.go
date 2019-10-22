package rest

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Routes assigns a path to a request handler.
func (api *API) Routes() *chi.Mux {
	// TODO @gritt: Test this!
	r := chi.NewRouter()

	// TODO @gritt: Create middleware to handle authenticated users

	r.Route("/v1", func(r chi.Router) {
		r.Method(http.MethodPost, "/transaction", api.HandleCreateTransaction())
		r.Method(http.MethodGet, "/transaction", api.HandleListTransaction())
	})

	return r
}
