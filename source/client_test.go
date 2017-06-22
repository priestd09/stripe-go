package source

import (
	"testing"

	assert "github.com/stretchr/testify/require"
	stripe "github.com/stripe/stripe-go"
	_ "github.com/stripe/stripe-go/testing"
)

func TestSourceGet(t *testing.T) {
	source, err := Get("gold", nil)
	assert.Nil(t, err)
	assert.NotNil(t, source)
}

func TestSourceNew(t *testing.T) {
	source, err := New(&stripe.SourceObjectParams{
		Type:     "bitcoin",
		Amount:   1000,
		Currency: "USD",
		Owner: &stripe.SourceOwnerParams{
			Email: "jenny.rosen@example.com",
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, source)
}

func TestSourceUpdate(t *testing.T) {
	source, err := Update("gold", &stripe.SourceObjectParams{
		Owner: &stripe.SourceOwnerParams{
			Email: "jenny.rosen@example.com",
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, source)
}
