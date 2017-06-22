package currency

import (
	"testing"

	assert "github.com/stretchr/testify/require"
	stripe "github.com/stripe/stripe-go"
	_ "github.com/stripe/stripe-go/testing"
)

func TestCurrencyWithChargeParams(t *testing.T) {
	params := &stripe.ChargeParams{
		Currency: USD,
	}
	assert.Equal(t, USD, params.Currency)
}
