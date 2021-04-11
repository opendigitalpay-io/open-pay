package trans

import (
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type Transfer struct {
	ID            uint64
	OrderID       uint64
	CustomerID    uint64
	SourceID      uint64
	DestinationID uint64
	Type          Type
	Amount        int64
	Currency      string
	Status        tcc.STATUS
	Metadata      map[string]interface{}
	CreatedAt     int64
	UpdatedAt     int64
}

type Type string

const (
	ORDER  Type = "ORDER"
	TOP_UP Type = "TOP_UP"
	REFUND Type = "REFUND"
)

var types = [...]string{
	"ORDER",
	"TOP_UP",
	"REFUND",
}

func (t *Type) String() string {
	x := string(*t)
	for _, v := range types {
		if v == x {
			return x
		}
	}
	return ""
}
