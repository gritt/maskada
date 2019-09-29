package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gritt/maskada/test"
)

func TestTransaction_validate(t *testing.T) {
	amount := test.RandomNumber()
	name := test.RandomName()
	invalidType := 1000

	tests := map[string]struct {
		givenTrs Transaction
		wantErr  error
	}{
		"should return err when missing amount": {
			givenTrs: Transaction{},
			wantErr:  errors.New("Transaction.Validate: missing amount"),
		},
		"should return err when zero amount": {
			givenTrs: Transaction{
				Amount: 0,
			},
			wantErr: errors.New("Transaction.Validate: missing amount"),
		},
		"should return err when negative amount": {
			givenTrs: Transaction{
				Amount: -1,
			},
			wantErr: errors.New("Transaction.Validate: missing amount"),
		},
		"should return err when missing type": {
			givenTrs: Transaction{
				Amount: amount,
			},
			wantErr: errors.New("Transaction.Validate: missing type"),
		},
		"should return err when invalid type": {
			givenTrs: Transaction{
				Amount: amount,
				Type:   invalidType,
			},
			wantErr: errors.New("Transaction.Validate: missing type"),
		},
		"should return err when invalid category": {
			givenTrs: Transaction{
				Amount: amount,
				Type:   Debit,
			},
			wantErr: errors.New("Transaction.Validate: missing category"),
		},
		"should return no err when valid debit transaction given": {
			givenTrs: Transaction{
				Amount:   amount,
				Type:     Debit,
				Category: Category{Name: name},
			},
		},
		"should return no err when valid credit transaction given": {
			givenTrs: Transaction{
				Amount:   amount,
				Type:     Credit,
				Category: Category{Name: name},
			},
		},
		"should return no err when valid income transaction given": {
			givenTrs: Transaction{
				Amount:   amount,
				Type:     Income,
				Category: Category{Name: name},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			gotErr := tc.givenTrs.Validate()

			if tc.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, tc.wantErr.Error())
			}
		})
	}
}
