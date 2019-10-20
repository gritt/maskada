package core

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gritt/maskada/test"
)

func TestCreateTransactionUseCase_CreateTransaction(t *testing.T) {
	amount := test.RandomNumber()
	name := test.RandomName()

	transaction := Transaction{
		Amount: amount,
		Type:   Debit,
		Category: Category{
			Name: name,
		},
	}

	wantTransaction := Transaction{
		ID:     test.RandomNumber(),
		Amount: amount,
		Type:   Debit,
		Category: Category{
			Name: name,
		},
		Date: time.Now(),
	}

	tests := map[string]func(t *testing.T, m *mockRepository){
		"should return err when invalid transaction": func(t *testing.T, m *mockRepository) {
			// arrange
			uc := NewCreateTransactionUseCase(m)

			// act
			got, gotErr := uc.Create(Transaction{})

			// assert
			assert.Empty(t, got)
			assert.EqualError(t, gotErr, "Create failed: Transaction.Validate: missing amount")
		},
		"should return err when fails to create transaction": func(t *testing.T, m *mockRepository) {
			// arrange
			m.On("Create", transaction).Return(Transaction{}, errors.New("Repository.Create: err"))
			uc := NewCreateTransactionUseCase(m)

			// act
			got, gotErr := uc.Create(transaction)

			// assert
			assert.Empty(t, got)
			assert.EqualError(t, gotErr, "Create failed: Repository.Create: err")
		},
		"should return transaction when succeed": func(t *testing.T, m *mockRepository) {
			// arrange
			m.On("Create", transaction).Return(wantTransaction, nil)
			uc := NewCreateTransactionUseCase(m)

			// act
			got, gotErr := uc.Create(transaction)

			// assert
			assert.NoError(t, gotErr)
			assert.Equal(t, wantTransaction, got)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			m := new(mockRepository)

			// act
			run(t, m)

			// assert
			m.AssertExpectations(t)
		})
	}
}

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Create(t Transaction) (Transaction, error) {
	args := m.Called(t)
	return args.Get(0).(Transaction), args.Error(1)
}
