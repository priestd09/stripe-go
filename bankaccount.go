package stripe

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/stripe/stripe-go/form"
)

// BankAccountStatus is the list of allowed values for the bank account's status.
// Allowed values are "new", "verified", "validated", "errored".
type BankAccountStatus string

// BankAccountParams is the set of parameters that can be used when updating a
// bank account.
//
// Note that while form annotations are used for updates, bank accounts have
// some unusual logic on creates that necessitates manual handling of all
// parameters. See AppendToAsSourceOrExternalAccount.
type BankAccountParams struct {
	Params `form:"*"`

	// The identifier of the parent account under which bank accounts are
	// nested.
	AccountID string `form:"-"`

	Customer string `form:"-"`

	// A token referencing an external account like one returned from
	// Stripe.js.
	Token string `form:"-"`

	// Information on an external account to reference. Only used if `Token`
	// is not provided.
	Account           string `form:"account_number"`
	AccountHolderName string `form:"account_holder_name"`
	AccountHolderType string `form:"account_holder_type"`
	Country           string `form:"country"`
	Currency          string `form:"currency"`
	Default           bool   `form:"default_for_currency"`
	Routing           string `form:"routing_number"`
}

// AppendToAsSourceOrExternalAccount appends the given BankAccountParams as
// either a source or external account.
//
// It may look like an AppendTo from the form package, but it's not, and is
// only used in the special case where we use `bankaccount.New`. It's needed
// because we have some weird encoding logic here that can't be handled by the
// form package (and it's special enough that it wouldn't be desirable to have
// it do so).
//
// This is not a pattern that we want to push forward, and this largely exists
// because the bank accounts endpoint is a little unusual. There is one other
// resource like it, which is cards.
func (a *BankAccountParams) AppendToAsSourceOrExternalAccount(body *form.Values) {
	isCustomer := len(a.Customer) > 0

	var sourceType string
	if isCustomer {
		sourceType = "source"
	} else {
		sourceType = "external_account"
	}

	// Use token (if exists) or a dictionary containing a user’s bank account details.
	if len(a.Token) > 0 {
		body.Add(sourceType, a.Token)

		if a.Default {
			body.Add("default_for_currency", strconv.FormatBool(a.Default))
		}
	} else {
		body.Add(sourceType+"[object]", "bank_account")
		body.Add(sourceType+"[country]", a.Country)
		body.Add(sourceType+"[account_number]", a.Account)
		body.Add(sourceType+"[currency]", a.Currency)

		if isCustomer {
			body.Add(sourceType+"[account_holder_name]", a.AccountHolderName)
			body.Add(sourceType+"[account_holder_type]", a.AccountHolderType)
		}

		if len(a.Routing) > 0 {
			body.Add(sourceType+"[routing_number]", a.Routing)
		}

		if a.Default {
			body.Add(sourceType+"[default_for_currency]", strconv.FormatBool(a.Default))
		}
	}
}

// BankAccountListParams is the set of parameters that can be used when listing bank accounts.
type BankAccountListParams struct {
	ListParams `form:"*"`

	// The identifier of the parent account under which the bank accounts are
	// nested. Either AccountID or Customer should be populated.
	AccountID string `form:"-"`

	// The identifier of the parent customer under which the bank accounts are
	// nested. Either AccountID or Customer should be populated.
	Customer string `form:"-"`
}

// BankAccount represents a Stripe bank account.
type BankAccount struct {
	ID                string            `json:"id"`
	Name              string            `json:"bank_name"`
	AccountHolderName string            `json:"account_holder_name"`
	AccountHolderType string            `json:"account_holder_type"`
	Country           string            `json:"country"`
	Currency          Currency          `json:"currency"`
	Default           bool              `json:"default_for_currency"`
	LastFour          string            `json:"last4"`
	Fingerprint       string            `json:"fingerprint"`
	Status            BankAccountStatus `json:"status"`
	Routing           string            `json:"routing_number"`
	Deleted           bool              `json:"deleted"`
	Customer          *Customer         `json:"customer"`
	Meta              map[string]string `json:"metadata"`
}

// BankAccountList is a list object for bank accounts.
type BankAccountList struct {
	ListMeta
	Values []*BankAccount `json:"data"`
}

// Display implements Displayer.Display.
func (b *BankAccount) Display() string {
	return fmt.Sprintf("Bank account ending in %s", b.LastFour)
}

// UnmarshalJSON handles deserialization of a BankAccount.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (b *BankAccount) UnmarshalJSON(data []byte) error {
	type bankAccount BankAccount
	var bb bankAccount
	err := json.Unmarshal(data, &bb)
	if err == nil {
		*b = BankAccount(bb)
	} else {
		// the id is surrounded by "\" characters, so strip them
		b.ID = string(data[1 : len(data)-1])
	}

	return nil
}
