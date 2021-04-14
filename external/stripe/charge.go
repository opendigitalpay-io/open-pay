package stripe

type StripeChargeRequest struct {
	Amount    int64
	Currency  string
	Source    string
	RequestID string
}

type StripeCaptureRequest struct {
	ChargeID  string
	RequestID string
}
