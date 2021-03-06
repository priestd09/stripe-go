package paymentsource

import (
	"testing"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/bankaccount"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/currency"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/source"
	"github.com/stripe/stripe-go/token"
	. "github.com/stripe/stripe-go/utils"
)

func init() {
	stripe.Key = GetTestKey()
}

func TestSourceCardNew(t *testing.T) {
	customerParams := &stripe.CustomerParams{}
	cust, err := customer.New(customerParams)

	if err != nil {
		t.Error(err)
	}

	sourceParams := &stripe.CustomerSourceParams{
		Customer: cust.ID,
	}
	sourceParams.SetSource("tok_visa")

	source, err := New(sourceParams)

	if err != nil {
		t.Error(err)
	}

	target := source.Card

	if target.LastFour != "4242" {
		t.Errorf("Unexpected last four %q for card number %v\n", target.LastFour, sourceParams.Source.Card.Number)
	}

	targetCust, err := customer.Get(cust.ID, nil)

	if err != nil {
		t.Error(err)
	}

	if targetCust.Sources.Count != 1 {
		t.Errorf("Unexpected number of sources %v\n", targetCust.Sources.Count)
	}

	customer.Del(cust.ID, nil)
}

func TestSourceBankAccountNew(t *testing.T) {
	baTok, err := token.New(&stripe.TokenParams{
		Bank: &stripe.BankAccountParams{
			Country:           "US",
			Currency:          "usd",
			Routing:           "110000000",
			Account:           "000123456789",
			AccountHolderName: "Jane Austen",
			AccountHolderType: "individual",
		},
	})

	if baTok.Bank.AccountHolderName != "Jane Austen" {
		t.Errorf("Bank account token name %q was not Jane Austen as expected.", baTok.Bank.AccountHolderName)
	}

	if err != nil {
		t.Error(err)
	}

	customerParams := &stripe.CustomerParams{}
	customer, err := customer.New(customerParams)
	if err != nil {
		t.Error(err)
	}

	source, err := New(&stripe.CustomerSourceParams{
		Customer: customer.ID,
		Source: &stripe.SourceParams{
			Token: baTok.ID,
		},
	})

	if source.BankAccount.AccountHolderName != "Jane Austen" {
		t.Errorf("Bank account name %q was not Jane Austen as expected.", source.BankAccount.Name)
	}
}

func TestSourceCardGet(t *testing.T) {
	customerParams := &stripe.CustomerParams{
		Email: "SomethingIdentifiable@gmail.om",
	}
	customerParams.SetSource("tok_visa")
	cust, err := customer.New(customerParams)

	if err != nil {
		t.Error(err)
	}

	source, err := Get(cust.DefaultSource.ID, &stripe.CustomerSourceParams{Customer: cust.ID})

	if err != nil {
		t.Error(err)
	}

	target := source.Card

	if target.LastFour != "4242" {
		t.Errorf("Unexpected last four %q for card number %v\n", target.LastFour, customerParams.Source.Card.Number)
	}

	if target.Brand != card.Visa {
		t.Errorf("Card brand %q does not match expected value\n", target.Brand)
	}

	customer.Del(cust.ID, nil)
}

func TestSourceBankAccountGet(t *testing.T) {
	baTok, err := token.New(&stripe.TokenParams{
		Bank: &stripe.BankAccountParams{
			Country:           "US",
			Currency:          "usd",
			Routing:           "110000000",
			Account:           "000123456789",
			AccountHolderName: "Jane Austen",
			AccountHolderType: "individual",
		},
	})

	if baTok.Bank.AccountHolderName != "Jane Austen" {
		t.Errorf("Bank account token name %q was not Jane Austen as expected.", baTok.Bank.AccountHolderName)
	}

	if err != nil {
		t.Error(err)
	}

	customerParams := &stripe.CustomerParams{}
	customer, err := customer.New(customerParams)
	if err != nil {
		t.Error(err)
	}

	src, err := New(&stripe.CustomerSourceParams{
		Customer: customer.ID,
		Source: &stripe.SourceParams{
			Token: baTok.ID,
		},
	})

	source, err := Get(src.ID, &stripe.CustomerSourceParams{Customer: customer.ID})

	if source.BankAccount.AccountHolderName != "Jane Austen" {
		t.Errorf("Bank account name %q was not Jane Austen as expected.", source.BankAccount.Name)
	}
}

func TestSourceCardDel(t *testing.T) {
	customerParams := &stripe.CustomerParams{}
	customerParams.SetSource("tok_visa")

	cust, _ := customer.New(customerParams)

	sourceDel, err := Del(cust.DefaultSource.ID, &stripe.CustomerSourceParams{Customer: cust.ID})
	if err != nil {
		t.Error(err)
	}

	if !sourceDel.Deleted {
		t.Errorf("Source id %q expected to be marked as deleted on the returned resource\n", sourceDel.ID)
	}

	customer.Del(cust.ID, nil)
}

func TestSourceBankAccountDel(t *testing.T) {
	baTok, err := token.New(&stripe.TokenParams{
		Bank: &stripe.BankAccountParams{
			Country:           "US",
			Currency:          "usd",
			Routing:           "110000000",
			Account:           "000123456789",
			AccountHolderName: "Jane Austen",
			AccountHolderType: "individual",
		},
	})

	if baTok.Bank.AccountHolderName != "Jane Austen" {
		t.Errorf("Bank account name %q was not Jane Austen as expected.", baTok.Bank.AccountHolderName)
	}

	if err != nil {
		t.Error(err)
	}

	customerParams := &stripe.CustomerParams{}
	customer, err := customer.New(customerParams)
	if err != nil {
		t.Error(err)
	}

	source, err := New(&stripe.CustomerSourceParams{
		Customer: customer.ID,
		Source: &stripe.SourceParams{
			Token: baTok.ID,
		},
	})

	sourceDel, err := Del(source.ID, &stripe.CustomerSourceParams{Customer: customer.ID})

	if !sourceDel.Deleted {
		t.Errorf("Source id %q expected to be marked as deleted on the returned resource\n", sourceDel.ID)
	}
}

func TestSourceCardUpdate(t *testing.T) {
	customerParams := &stripe.CustomerParams{}
	customerParams.SetSource("tok_visa")

	cust, err := customer.New(customerParams)

	if err != nil {
		t.Error(err)
	}

	sourceParams := &stripe.CustomerSourceParams{
		Customer: cust.ID,
	}
	sourceParams.SetSource(&stripe.CardParams{
		Name: "Updated Name",
	})

	source, err := Update(cust.DefaultSource.ID, sourceParams)

	if err != nil {
		t.Error(err)
		return
	}

	target := source.Card

	if target.Name != sourceParams.Source.Card.Name {
		t.Errorf("Card name %q does not match expected name %q\n", target.Name, sourceParams.Source.Card.Name)
	}

	customer.Del(cust.ID, nil)
}

func TestSourceBankAccountVerify(t *testing.T) {
	baTok, err := token.New(&stripe.TokenParams{
		Bank: &stripe.BankAccountParams{
			Country:           "US",
			Currency:          "usd",
			Routing:           "110000000",
			Account:           "000123456789",
			AccountHolderName: "Jane Austen",
			AccountHolderType: "individual",
		},
	})

	if baTok.Bank.AccountHolderName != "Jane Austen" {
		t.Errorf("Bank account name %q was not Jane Austen as expected.", baTok.Bank.AccountHolderName)
	}

	if err != nil {
		t.Error(err)
	}

	customerParams := &stripe.CustomerParams{}
	cust, err := customer.New(customerParams)
	if err != nil {
		t.Error(err)
	}

	source, err := New(&stripe.CustomerSourceParams{
		Customer: cust.ID,
		Source: &stripe.SourceParams{
			Token: baTok.ID,
		},
	})

	amounts := [2]uint8{32, 45}

	verifyParams := &stripe.SourceVerifyParams{
		Customer: cust.ID,
		Amounts:  amounts,
	}

	sourceVerified, err := Verify(source.ID, verifyParams)

	if err != nil {
		t.Error(err)
		return
	}

	target := sourceVerified.BankAccount

	if target.Status != bankaccount.VerifiedAccount {
		t.Errorf("Status (%q) does not match expected (%q) ", target.Status, bankaccount.VerifiedAccount)
	}

	customer.Del(cust.ID, nil)
}

func TestSourceList(t *testing.T) {
	customerParams := &stripe.CustomerParams{}
	customerParams.SetSource("tok_amex")

	cust, _ := customer.New(customerParams)

	sourceParams := &stripe.CustomerSourceParams{
		Customer: cust.ID,
	}
	sourceParams.SetSource("tok_visa")

	New(sourceParams)

	i := List(&stripe.SourceListParams{Customer: cust.ID})
	for i.Next() {
		paymentSource := i.PaymentSource()

		if paymentSource == nil {
			t.Error("No nil values expected")
		}

		if paymentSource.Card == nil {
			t.Error("No nil values expected")
		}

		if i.Meta() == nil {
			t.Error("No metadata returned")
		}
	}
	if err := i.Err(); err != nil {
		t.Error(err)
	}

	customer.Del(cust.ID, nil)
}

func TestSourceObjectNewGet(t *testing.T) {
	sourceParams := &stripe.SourceObjectParams{
		Type:     "bitcoin",
		Amount:   1000,
		Currency: currency.USD,
		Owner: &stripe.SourceOwnerParams{
			Email: "do+fill_now@stripe.com",
		},
	}

	s, err := source.New(sourceParams)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	customerParams := &stripe.CustomerParams{}
	customer, err := customer.New(customerParams)
	if err != nil {
		t.Error(err)
	}

	src, err := New(&stripe.CustomerSourceParams{
		Customer: customer.ID,
		Source: &stripe.SourceParams{
			Token: s.ID,
		},
	})

	if src.SourceObject.Owner.Email != "do+fill_now@stripe.com" {
		t.Errorf("Source object owner email %q was not as expected.",
			src.SourceObject.Owner.Email)
	}

	src, err = Get(src.ID, &stripe.CustomerSourceParams{
		Customer: customer.ID,
	})
	if err != nil {
		t.Error(err)
	}

	if src.SourceObject.Owner.Email != "do+fill_now@stripe.com" {
		t.Errorf("Source object owner email %q was not as expected.",
			src.SourceObject.Owner.Email)
	}
}
