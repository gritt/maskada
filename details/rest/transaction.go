package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gritt/maskada/core"
)

type skeleton struct {
	ID       int       `json:"id"`
	Amount   int       `json:"amount"`
	Type     int       `json:"type"`
	Category string    `json:"category"`
	Date     time.Time `json:"date"`
	Name     string    `json:"name"`
}

// HandleCreateTransaction receives the request and call the use case to create a transaction.
func (api *API) HandleCreateTransaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil || r.Method != http.MethodPost {
			respond(w, `{"error": "HandleCreateTransaction failed: invalid request"}`, http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			respond(w, `{"error": "HandleCreateTransaction failed: could not read body"}`, http.StatusBadRequest)
			return
		}

		payload := skeleton{}
		err = json.Unmarshal(body, &payload)
		if err != nil {
			respond(w, `{"error": "HandleCreateTransaction failed: could not decode payload"}`, http.StatusBadRequest)
			return
		}

		trs, err := api.TransactionCreator.Create(core.Transaction{
			Amount:   payload.Amount,
			Type:     payload.Type,
			Category: core.Category{Name: payload.Category},
			Date:     payload.Date,
			Name:     payload.Name,
		})
		if err != nil {
			respond(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		res := skeleton{
			ID:       trs.ID,
			Amount:   trs.Amount,
			Type:     trs.Type,
			Category: trs.Category.Name,
			Date:     trs.Date,
			Name:     trs.Name,
		}
		jsonRes, _ := json.Marshal(&res)
		respond(w, string(jsonRes), http.StatusCreated)
	}
}

// HandleListTransaction receives the request and call the use case to list transactions.
func (api *API) HandleListTransaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respond(w, `{"error": "HandleCreateTransaction failed: invalid request"}`, http.StatusBadRequest)
			return
		}

		trsl, err := api.TransactionLister.List()
		if err != nil {
			respond(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		res := []skeleton{}
		for _, trs := range trsl {
			res = append(res, skeleton{
				ID:       trs.ID,
				Amount:   trs.Amount,
				Type:     trs.Type,
				Category: trs.Category.Name,
				Date:     trs.Date,
				Name:     trs.Name,
			})
		}
		jsonRes, _ := json.Marshal(&res)
		respond(w, string(jsonRes), http.StatusOK)
	}
}

func respond(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(msg))
}
