package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gritt/maskada/core"
)

func TestAPI_handleCreateTransaction(t *testing.T) {
	payload := `{"amount": 32, "type": 2, "category": "Food"}`

	transaction := core.Transaction{
		Amount:   32,
		Type:     core.Credit,
		Category: core.Category{Name: "Food"},
		Date:     time.Time{},
	}

	createdTransaction := core.Transaction{
		ID:       1,
		Amount:   transaction.Amount,
		Type:     transaction.Type,
		Category: transaction.Category,
		Date:     time.Now().UTC(),
		Name:     "",
	}

	tests := map[string]func(*testing.T){
		"should return error when invalid request": func(t *testing.T) {
			// arrange
			m := new(mockTransactionCreator)
			api := NewAPI(m)

			rr := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", nil)

			// act
			api.HandleCreateTransaction()(rr, r)

			// assert
			assert.Equal(t, http.StatusBadRequest, rr.Code)
			assert.Equal(t, `{"error": "HandleCreateTransaction failed: invalid request"}`, rr.Body.String())
			assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
			m.AssertExpectations(t)
		},
		"should return error when invalid body": func(t *testing.T) {
			// arrange
			m := new(mockTransactionCreator)
			api := NewAPI(m)

			rr := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{{}`))

			// act
			api.HandleCreateTransaction()(rr, r)

			// assert
			assert.Equal(t, http.StatusBadRequest, rr.Code)
			assert.Equal(t, `{"error": "HandleCreateTransaction failed: could not decode payload"}`, rr.Body.String())
			assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
			m.AssertExpectations(t)
		},
		"should return error when create returns error": func(t *testing.T) {
			// arrange
			m := new(mockTransactionCreator)
			api := NewAPI(m)
			m.On("Create", transaction).Return(core.Transaction{}, errors.New("Create failed: err"))

			rr := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(payload))

			// act
			api.HandleCreateTransaction()(rr, r)

			// assert
			assert.Equal(t, http.StatusInternalServerError, rr.Code)
			assert.Equal(t, `{"error": "Create failed: err"}`, rr.Body.String())
			assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
			m.AssertExpectations(t)
		},
		"should return success": func(t *testing.T) {
			// arrange
			m := new(mockTransactionCreator)
			api := NewAPI(m)
			m.On("Create", transaction).Return(createdTransaction, nil)

			rr := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(payload))

			// act
			api.HandleCreateTransaction()(rr, r)

			wantBody, _ := json.Marshal(&skeleton{
				ID:       createdTransaction.ID,
				Amount:   createdTransaction.Amount,
				Type:     createdTransaction.Type,
				Category: createdTransaction.Category.Name,
				Date:     createdTransaction.Date,
				Name:     createdTransaction.Name,
			})

			// assert
			assert.Equal(t, http.StatusCreated, rr.Code)
			assert.Equal(t, string(wantBody), rr.Body.String())
			assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
			m.AssertExpectations(t)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}
