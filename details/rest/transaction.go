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

// HandleCreateTransaction receive the request and call the use case to create a transaction.
func (api *API) HandleCreateTransaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond := func(message string, status int) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(status)
			_, _ = w.Write([]byte(message))
		}

		if r.Body == nil {
			respond(`{"error": "HandleCreateTransaction failed: invalid request"}`, http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			respond(`{"error": "HandleCreateTransaction failed: could not read body"}`, http.StatusBadRequest)
			return
		}

		trs := skeleton{}
		err = json.Unmarshal(body, &trs)
		if err != nil {
			respond(`{"error": "HandleCreateTransaction failed: could not decode payload"}`, http.StatusBadRequest)
			return
		}

		transaction, err := api.TransactionCreator.Create(core.Transaction{
			Amount:   trs.Amount,
			Type:     trs.Type,
			Category: core.Category{Name: trs.Category},
			Date:     trs.Date,
			Name:     trs.Name,
		})
		if err != nil {
			respond(fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		response, _ := json.Marshal(&skeleton{
			ID:       transaction.ID,
			Amount:   transaction.Amount,
			Type:     transaction.Type,
			Category: transaction.Category.Name,
			Date:     transaction.Date,
			Name:     transaction.Name,
		})
		respond(string(response), http.StatusCreated)
	}
}
