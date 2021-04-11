package domain

import "github.com/opendigitalpay-io/open-pay/internal/tcc"

type Refund struct {
	ID          uint64
	OrderID     uint64
	Amount      int64
	Status      tcc.STATUS
	RefundCount int32
	Metadata    map[string]interface{}
	CreatedAt   int64
	UpdatedAt   int64
}
