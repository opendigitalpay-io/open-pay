package stripe

import (
	"fmt"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/charge"
)

type Adapter interface {
	Authorize(StripeChargeRequest) (string, error)
	Capture(StripeCaptureRequest) error
}

type adapter struct {
}

func NewAdapter() Adapter {
	stripe.Key = "sk_test_PFluMxE3L4WDrISWmkDO3n4A"
	return &adapter{}
}

func (a *adapter) Authorize(request StripeChargeRequest) (string, error) {
	params := &stripe.ChargeParams{
		Amount:   stripe.Int64(request.Amount),
		Currency: stripe.String(request.Currency),
		Source: &stripe.SourceParams{
			Token: stripe.String("tok_visa"),
		},
	}
	params.SetIdempotencyKey(request.RequestID)
	c, err := charge.New(params)
	if err != nil {
		fmt.Printf("Other error occurred: %v\n", err.Error())
		// FIXME: error handling logic
		// Try to safely cast a generic error to a stripe.Error so that we can get at
		// some additional Stripe-specific information about what went wrong.
		//if stripeErr, ok := err.(*stripe.Error); ok {
		//	// The Code field will contain a basic identifier for the failure.
		//	switch stripeErr.Code {
		//	case stripe.ErrorCodeCardDeclined:
		//	case stripe.ErrorCodeExpiredCard:
		//	case stripe.ErrorCodeIncorrectCVC:
		//	case stripe.ErrorCodeIncorrectZip:
		//		// etc.
		//	}
		//
		//	// The Err field can be coerced to a more specific error type with a type
		//	// assertion. This technique can be used to get more specialized
		//	// information for certain errors.
		//	if cardErr, ok := stripeErr.Err.(*stripe.CardError); ok {
		//		fmt.Printf("Card was declined with code: %v\n", cardErr.DeclineCode)
		//	} else {
		//		fmt.Printf("Other Stripe error occurred: %v\n", stripeErr.Error())
		//	}
		//} else {
		//	fmt.Printf("Other error occurred: %v\n", err.Error())
		//}
		return "", err
	}

	return c.ID, nil
}

func (a *adapter) Capture(request StripeCaptureRequest) error {
	//TODO: need testing not sure if empty param here works
	params := &stripe.CaptureParams{}
	params.SetIdempotencyKey(request.RequestID)
	_, err := charge.Capture(request.ChargeID, params)

	if err != nil {
		fmt.Printf("Other error occurred: %v\n", err.Error())
		return err
	}

	return nil
}
