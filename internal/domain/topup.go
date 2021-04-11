package domain

import "github.com/opendigitalpay-io/open-pay/internal/tcc"

type TopUp struct {
	ID              uint64
	CustomerID      uint64
	PaymentMethodID uint64
	Amount          int64
	Currency        string
	Status          tcc.STATUS
	Metadata        map[string]interface{}
	CreatedAt       int64
	UpdatedAt       int64
}
