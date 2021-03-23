package api

type AddTopupUriParameter struct {
	ID uint64 `uri:"id" binding:"required"`
}

type AddTopupRequest struct {
	Amount          int64  `json:"amount" binding:"required"`
	Currency        string `json:"currency" binding:"required"`
	PaymentMethodID uint64 `json:"paymentMethodId"`
}

type AddTopupResponse struct {
	ID uint64 `json:"id"`
}
