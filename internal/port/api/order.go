package api

type AddOrderRequest struct {
	CustomerID    uint64                 `json:"customerId" binding:"required"`
	BusinessID    uint64                 `json:"businessId" binding:"required"`
	Amount        int64                  `json:"amount" binding:"required"`
	Currency      string                 `json:"currency" binding:"required"`
	ReferenceID   string                 `json:"referenceId" binding:"required"`
	CustomerEmail string                 `json:"customerEmail" binding:"required"`
	LineItems     []LineItem             `json:"lineItems"`
	Metadata      map[string]interface{} `json:"metadata"`
}

type GetOrderURIParameter struct {
	ID uint64 `uri:"id" binding:"required"`
}

type GerOrderResponse struct {
	ID            uint64                 `json:"id"`
	CustomerID    uint64                 `json:"customerId"`
	BusinessID    uint64                 `json:"businessId"`
	Amount        int64                  `json:"amount"`
	Currency      string                 `json:"currency"`
	ReferenceID   string                 `json:"referenceId"`
	CustomerEmail string                 `json:"customerEmail"`
	LineItems     []LineItem             `json:"lineItems"`
	Metadata      map[string]interface{} `json:"metadata"`
	CreatedAt     int64                  `json:"createdAt"`
}

type LineItem struct {
	Name           string         `json:"name"`
	Quantity       int32          `json:"quantity"`
	BasePriceMoney BasePriceMoney `json:"basePriceMoney"`
}

type BasePriceMoney struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type AddOrderResponse struct {
	OrderID uint64 `json:"id"`
}
