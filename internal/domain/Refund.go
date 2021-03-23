package domain

type Refund struct {
	ID          uint64
	OrderID     uint64
	Amount      int64
	Status      STATUS
	RefundCount int32
	Metadata    []byte
	CreatedAt   int64
	UpdatedAt   int64
}