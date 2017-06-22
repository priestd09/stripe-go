package bankaccount

import (
	"testing"

	assert "github.com/stretchr/testify/require"
	stripe "github.com/stripe/stripe-go"
	_ "github.com/stripe/stripe-go/testing"
)

func TestBankAccountDel(t *testing.T) {
	account, err := Del("ba_123", &stripe.BankAccountParams{
		AccountID: "acct_123",
	})
	assert.Nil(t, err)
	assert.NotNil(t, account)
}

func TestBankAccountGet(t *testing.T) {
	account, err := Get("ba_123", &stripe.BankAccountParams{AccountID: "acct_123"})
	assert.Nil(t, err)
	assert.NotNil(t, account)
}

func TestBankAccountListByCustomer(t *testing.T) {
	i := List(&stripe.BankAccountListParams{Customer: "cus_123"})

	// Verify that we can get at least one plan
	assert.True(t, i.Next())
	assert.Nil(t, i.Err())
	assert.NotNil(t, i.BankAccount())
}

func TestBankAccountNew(t *testing.T) {
	account, err := New(&stripe.BankAccountParams{
		Customer: "cus_123",
		Default:  true,
		Token:    "tok_123",
	})
	assert.Nil(t, err)
	assert.NotNil(t, account)
}

func TestBankAccountUpdate(t *testing.T) {
	account, err := Update("ba_123", &stripe.BankAccountParams{
		AccountID: "acct_123",
		Default:   true,
	})
	assert.Nil(t, err)
	assert.NotNil(t, account)
}
