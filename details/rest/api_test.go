package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gritt/maskada/core"
)

func TestNewAPI(t *testing.T) {
	// arrange
	c := new(mockTransactionCreator)
	l := new(mockTransactionLister)

	// act
	got := NewAPI(c, l)

	want := &API{
		TransactionCreator: c,
		TransactionLister:  l,
	}

	// assert
	assert.Equal(t, want, got)
}

type mockTransactionCreator struct {
	mock.Mock
}

func (m *mockTransactionCreator) Create(t core.Transaction) (core.Transaction, error) {
	args := m.Called(t)
	return args.Get(0).(core.Transaction), args.Error(1)
}

type mockTransactionLister struct {
	mock.Mock
}

func (m *mockTransactionLister) List() ([]core.Transaction, error) {
	args := m.Called()
	return args.Get(0).([]core.Transaction), args.Error(1)
}
