package core

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gritt/maskada/test"
)

func TestTransaction_validate(t *testing.T) {
	amount := test.RandomNumber()
	name := test.RandomName()
	invalidType := 1000

	tests := map[string]func(t *testing.T){
		"when missing amount": func(t *testing.T) {
			// arrange
			trs := Transaction{}

			// act
			gotErr := trs.Validate()

			// assert
			assert.EqualError(t, gotErr, "Transaction.Validate: invalid amount")
		},
		"when zero amount": func(t *testing.T) {
			// arrange
			trs := Transaction{Amount: 0}

			// act
			gotErr := trs.Validate()

			// assert
			assert.EqualError(t, gotErr, "Transaction.Validate: invalid amount")
		},
		"when negative amount": func(t *testing.T) {
			// arrange
			trs := Transaction{Amount: -1}

			// act
			gotErr := trs.Validate()

			// assert
			assert.EqualError(t, gotErr, "Transaction.Validate: invalid amount")
		},
		"when missing type": func(t *testing.T) {
			// arrange
			trs := Transaction{
				Amount: amount,
			}

			// act
			gotErr := trs.Validate()

			// assert
			assert.EqualError(t, gotErr, "Transaction.Validate: invalid type")
		},
		"when invalid type": func(t *testing.T) {
			// arrange
			trs := Transaction{Amount: amount, Type: invalidType}

			// act
			gotErr := trs.Validate()

			// assert
			assert.EqualError(t, gotErr, "Transaction.Validate: invalid type")
		},
		"when invalid category": func(t *testing.T) {
			// arrange
			trs := Transaction{Amount: amount, Type: Debit}

			// act
			gotErr := trs.Validate()

			// assert
			assert.EqualError(t, gotErr, "Transaction.Validate: invalid category")
		},
		"when valid debit transaction given": func(t *testing.T) {
			// arrange
			trs := Transaction{Amount: amount, Type: Debit, Category: Category{Name: name}}

			// act / assert
			assert.NoError(t, trs.Validate())
		},
		"when valid credit transaction given": func(t *testing.T) {
			// arrange
			trs := Transaction{Amount: amount, Type: Credit, Category: Category{Name: name}}

			// act / assert
			assert.NoError(t, trs.Validate())
		},
		"when valid income transaction given": func(t *testing.T) {
			// arrange
			trs := Transaction{Amount: amount, Type: Income, Category: Category{Name: name}}

			// act / assert
			assert.NoError(t, trs.Validate())
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}
