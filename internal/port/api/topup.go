package api

type AddTopUpUriParameter struct {
	ID uint64 `uri:"id" binding:"required"`
}

type AddTopUpRequest struct {
	Amount          int64  `json:"amount" binding:"required"`
	Currency        string `json:"currency" binding:"required"`
	PaymentMethodID uint64 `json:"paymentMethodId"`
}

type AddTopUpResponse struct {
	ID uint64 `json:"id"`
}
