package stripe

import (
	"testing"

	assert "github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go/form"
)

func TestSourceParams_AppendTo(t *testing.T) {
	{
		params := &SourceParams{Token: "src_123"}
		body := &form.Values{}
		form.AppendTo(body, params)
		t.Logf("body = %+v", body)
		assert.Equal(t, []string{"src_123"}, body.Get("source"))
	}

	{
		params := &SourceParams{Card: &CardParams{Name: "a card"}}
		body := &form.Values{}
		form.AppendTo(body, params)
		t.Logf("body = %+v", body)
		assert.Equal(t, []string{"a card"}, body.Get("source[name]"))
		assert.Equal(t, []string{"card"}, body.Get("source[object]"))
	}
}
