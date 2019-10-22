package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gritt/maskada/core"
	"github.com/gritt/maskada/test"
)

var (
	testAmount   = test.RandomNumber()
	testCategory = test.RandomName()
	testName     = test.RandomName()

	testPayload = fmt.Sprintf(`{"amount": %d, "type": 2, "category": "%s", "name": "%s"}`, testAmount, testCategory, testName)

	testTrs = core.Transaction{
		Amount:   testAmount,
		Type:     core.Credit,
		Category: core.Category{Name: testCategory},
		Date:     time.Time{},
		Name:     testName,
	}

	testCreatedTrs = core.Transaction{
		ID:       test.RandomNumber(),
		Amount:   testTrs.Amount,
		Type:     testTrs.Type,
		Category: testTrs.Category,
		Date:     time.Now().UTC(),
		Name:     testName,
	}

	testCreatedTrsList = []core.Transaction{{
		ID:       testCreatedTrs.ID,
		Amount:   testCreatedTrs.Amount,
		Type:     testCreatedTrs.Type,
		Category: testCreatedTrs.Category,
		Date:     testCreatedTrs.Date,
		Name:     testCreatedTrs.Name,
	}}
)

func TestAPI_handleCreateTransaction(t *testing.T) {
	tests := map[string]func(*testing.T){
		"when invalid method": func(t *testing.T) {
			// arrange
			c := new(mockTransactionCreator)
			api := NewAPI(c, nil)

			rr := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(`{}`))

			// act
			api.HandleCreateTransaction()(rr, r)

			// assert
			assert.Equal(t, http.StatusBadRequest, rr.Code)
			assert.Equal(t, `{"error": "HandleCreateTransaction failed: invalid request"}`, rr.Body.String())
			assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
			c.AssertExpectations(t)
		},
		"when invalid request": func(t *testing.T) {
			// arrange
			c := new(mockTransactionCreator)
			api := NewAPI(c, nil)

			rr := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", nil)

			// act
			api.HandleCreateTransaction()(rr, r)

			// assert
			assert.Equal(t, http.StatusBadRequest, rr.Code)
			assert.Equal(t, `{"error": "HandleCreateTransaction failed: invalid request"}`, rr.Body.String())
			assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
			c.AssertExpectations(t)
		},
		"when invalid body": func(t *testing.T) {
			// arrange
			c := new(mockTransactionCreator)
			api := NewAPI(c, nil)

			rr := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{{}`))

			// act
			api.HandleCreateTransaction()(rr, r)

			// assert
			assert.Equal(t, http.StatusBadRequest, rr.Code)
			assert.Equal(t, `{"error": "HandleCreateTransaction failed: could not decode payload"}`, rr.Body.String())
			assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
			c.AssertExpectations(t)
		},
		"when create returns error": func(t *testing.T) {
			// arrange
			c := new(mockTransactionCreator)
			c.On("Create", testTrs).Return(core.Transaction{}, errors.New("Create failed: err"))
			api := NewAPI(c, nil)

			rr := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(testPayload))

			// act
			api.HandleCreateTransaction()(rr, r)

			// assert
			assert.Equal(t, http.StatusInternalServerError, rr.Code)
			assert.Equal(t, `{"error": "Create failed: err"}`, rr.Body.String())
			assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
			c.AssertExpectations(t)
		},
		"when succeed creating transaction": func(t *testing.T) {
			// arrange
			c := new(mockTransactionCreator)
			c.On("Create", testTrs).Return(testCreatedTrs, nil)
			api := NewAPI(c, nil)

			rr := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(testPayload))

			// act
			api.HandleCreateTransaction()(rr, r)

			want, _ := json.Marshal(&skeleton{
				ID:       testCreatedTrs.ID,
				Amount:   testCreatedTrs.Amount,
				Type:     testCreatedTrs.Type,
				Category: testCreatedTrs.Category.Name,
				Date:     testCreatedTrs.Date,
				Name:     testCreatedTrs.Name,
			})

			// assert
			assert.Equal(t, http.StatusCreated, rr.Code)
			assert.Equal(t, string(want), rr.Body.String())
			assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
			c.AssertExpectations(t)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func TestAPI_HandleListTransaction(t *testing.T) {
	tests := map[string]func(t *testing.T){
		"when invalid method": func(t *testing.T) {
			// arrange
			l := new(mockTransactionLister)
			api := NewAPI(nil, l)

			rr := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", nil)

			// act
			api.HandleListTransaction()(rr, r)

			// assert
			assert.Equal(t, http.StatusBadRequest, rr.Code)
		},
		"when list returns error": func(t *testing.T) {
			// arrange
			l := new(mockTransactionLister)
			l.On("List").Return([]core.Transaction{}, errors.New("List failed: err"))
			api := NewAPI(nil, l)

			rr := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "/", nil)

			// act
			api.HandleListTransaction()(rr, r)

			// assert
			assert.Equal(t, http.StatusInternalServerError, rr.Code)
			assert.Equal(t, `{"error": "List failed: err"}`, rr.Body.String())
			assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
			l.AssertExpectations(t)
		},
		"when succeed with empty list": func(t *testing.T) {
			// arrange
			l := new(mockTransactionLister)
			l.On("List").Return([]core.Transaction{}, nil)
			api := NewAPI(nil, l)

			rr := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "/", nil)

			// act
			api.HandleListTransaction()(rr, r)

			// assert
			assert.Equal(t, http.StatusOK, rr.Code)
			assert.Equal(t, `[]`, rr.Body.String())
			l.AssertExpectations(t)
		},
		"when succeed with list": func(t *testing.T) {
			// arrange
			l := new(mockTransactionLister)
			l.On("List").Return(testCreatedTrsList, nil)
			api := NewAPI(nil, l)

			rr := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "/", nil)

			// act
			api.HandleListTransaction()(rr, r)

			want, _ := json.Marshal(&[]skeleton{{
				ID:       testCreatedTrsList[0].ID,
				Amount:   testCreatedTrsList[0].Amount,
				Type:     testCreatedTrsList[0].Type,
				Category: testCreatedTrsList[0].Category.Name,
				Date:     testCreatedTrsList[0].Date,
				Name:     testCreatedTrsList[0].Name,
			}})

			// assert
			assert.Equal(t, http.StatusOK, rr.Code)
			assert.Equal(t, string(want), rr.Body.String())
			assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
			l.AssertExpectations(t)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}
