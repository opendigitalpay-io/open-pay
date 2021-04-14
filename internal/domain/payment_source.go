package domain

type PaymentSourceType string

const (
	TOKEN           PaymentSourceType = "TOKEN"
	PAYMENT_METHOD  PaymentSourceType = "PAYMENT_METHOD"
	BALANCE_ACCOUNT PaymentSourceType = "BALANCE_ACCOUNT"
	INTERACT        PaymentSourceType = "INTERACT"
)

type PaymentSource struct {
	Type PaymentSourceType
	ID   string
}
