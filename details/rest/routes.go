package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

// Routes assigns a path to a request handler.
func (api *API) Routes() *chi.Mux {
	// TODO @gritt: Test this!
	r := chi.NewRouter()

	mw := []func(http.Handler) http.Handler{}

	mw = append(
		mw,
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST"},
			AllowCredentials: true,
		}).Handler)

	r.Use(mw...)

	r.Route("/v1", func(r chi.Router) {
		r.Method(http.MethodPost, "/transaction", api.HandleCreateTransaction())
		r.Method(http.MethodGet, "/transaction", api.HandleListTransaction())
	})

	return r
}
