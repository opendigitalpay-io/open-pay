package domain

type Topup struct {
	ID              uint64
	CustomerID      uint64
	PaymentMethodID uint64
	Amount          int64
	Currency        string
	Status          STATUS
	Metadata        []byte
	CreatedAt       int64
	UpdatedAt       int64
}
