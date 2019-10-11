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
	type given struct {
		trs     Transaction
		arrange func(m *mockRepository)
	}

	type want struct {
		trs Transaction
		err error
	}

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

	tests := map[string]struct {
		given
		want
	}{
		"should return err when invalid transaction": {
			given: given{
				trs:     Transaction{},
				arrange: func(m *mockRepository) {},
			},
			want: want{
				err: errors.New("Create failed: Transaction.Validate: missing amount"),
			},
		},
		"should return err when fails to create transaction": {
			given: given{
				trs: transaction,
				arrange: func(m *mockRepository) {
					m.On("Create", transaction).Return(Transaction{}, errors.New("Repository.Create: err"))
				},
			},
			want: want{
				err: errors.New("Create failed: Repository.Create: err"),
			},
		},
		"should return transaction when succeed": {
			given: given{
				trs: transaction,
				arrange: func(m *mockRepository) {
					m.On("Create", transaction).Return(wantTransaction, nil)
				},
			},
			want: want{
				trs: wantTransaction,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			mockRepository := new(mockRepository)
			tc.given.arrange(mockRepository)

			uc := CreateTransactionUseCase{repository: mockRepository}

			// act
			gotTrs, gotErr := uc.Create(tc.given.trs)

			// assert
			if tc.want.err == nil {
				assert.NoError(t, gotErr)
				assert.Equal(t, tc.want.trs, gotTrs)
			} else {
				assert.EqualError(t, gotErr, tc.want.err.Error())
			}

			mockRepository.AssertExpectations(t)
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
