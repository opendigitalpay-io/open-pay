package api

type AddRefundUriParameter struct {
	OrderID uint64 `uri:"orderId" binding:"required"`
}

type AddRefundResponse struct {
	ID uint64 `json:"id"`
}
