package domain

type TopUp struct {
	ID              uint64
	CustomerID      uint64
	PaymentMethodID uint64
	Amount          int64
	Currency        string
	Status          STATUS
	Metadata        map[string]interface{}
	CreatedAt       int64
	UpdatedAt       int64
}
