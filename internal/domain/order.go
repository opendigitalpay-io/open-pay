package domain

import "github.com/opendigitalpay-io/open-pay/internal/tcc"

type OrderMode string

const (
	DIRECT   OrderMode = "DIRECT"
	INDIRECT OrderMode = "INDIRECT"
)

var orderModes = [...]string{
	"DIRECT",
	"INDIRECT",
}

func (s *OrderMode) String() string {
	x := string(*s)
	for _, v := range orderModes {
		if v == x {
			return x
		}
	}
	return ""
}

type Order struct {
	ID            uint64
	CustomerID    uint64
	MerchantID    uint64
	Amount        int64
	Currency      string
	ReferenceID   string
	CustomerEmail string
	Status        tcc.STATUS
	Mode          OrderMode
	Metadata      map[string]interface{}
	CreatedAt     int64
	UpdatedAt     int64
}

type LineItem struct {
	Name           string
	Quantity       int32
	BasePriceMoney BasePriceMoney
}

type BasePriceMoney struct {
	Amount   int64
	Currency string
}
