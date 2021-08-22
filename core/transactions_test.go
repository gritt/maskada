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
		"when invalid transaction": func(t *testing.T, m *mockRepository) {
			// arrange
			uc := NewCreateTransactionUseCase(m)

			// act
			got, gotErr := uc.Create(Transaction{})

			// assert
			assert.EqualError(t, gotErr, "Create failed: Transaction.Validate: invalid amount")
			assert.Empty(t, got)
		},
		"when repository fails to create transaction": func(t *testing.T, m *mockRepository) {
			// arrange
			m.On("Create", transaction).Return(Transaction{}, errors.New("Repository.Create: err"))
			uc := NewCreateTransactionUseCase(m)

			// act
			got, gotErr := uc.Create(transaction)

			// assert
			assert.EqualError(t, gotErr, "Create failed: Repository.Create: err")
			assert.Empty(t, got)
		},
		"when repository creates transaction": func(t *testing.T, m *mockRepository) {
			// arrange
			m.On("Create", transaction).Return(wantTransaction, nil)
			uc := NewCreateTransactionUseCase(m)

			// act
			got, gotErr := uc.Create(transaction)

			// assert
			assert.Equal(t, wantTransaction, got)
			assert.NoError(t, gotErr)
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

func TestListTransactionUseCase_List(t *testing.T) {
	tests := map[string]func(t *testing.T, m *mockRepository){
		"when repository fails to list transactions": func(t *testing.T, m *mockRepository) {
			// arrange
			m.On("Find").Return([]Transaction{}, errors.New("Repository.Find: err"))
			uc := NewListTransactionUseCase(m)

			// act
			got, gotErr := uc.List()

			// assert
			assert.EqualError(t, gotErr, "List failed: Repository.Find: err")
			assert.Empty(t, got)
		},
		"when repository returns empty list of transactions": func(t *testing.T, m *mockRepository) {
			// arrange
			m.On("Find").Return([]Transaction{}, nil)
			uc := NewListTransactionUseCase(m)

			// act
			got, gotErr := uc.List()

			// assert
			assert.Empty(t, got)
			assert.NoError(t, gotErr)
		},
		"when repository returns transactions": func(t *testing.T, m *mockRepository) {
			// arrange
			transaction1 := Transaction{
				ID:       test.RandomNumber(),
				Amount:   test.RandomNumber(),
				Type:     Credit,
				Category: Category{Name: test.RandomName()},
				Date:     time.Now(),
			}
			transaction2 := Transaction{
				ID:       test.RandomNumber(),
				Amount:   test.RandomNumber(),
				Type:     Credit,
				Category: Category{Name: test.RandomName()},
				Date:     time.Now().Add(time.Duration(1)),
			}

			m.On("Find").Return([]Transaction{transaction1, transaction2}, nil)
			uc := NewListTransactionUseCase(m)

			// act
			got, gotErr := uc.List()

			// assert
			assert.ElementsMatch(t, got, []Transaction{transaction1, transaction2})
			assert.NoError(t, gotErr)
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

func (m *mockRepository) Find() ([]Transaction, error) {
	args := m.Called()
	return args.Get(0).([]Transaction), args.Error(1)
}
