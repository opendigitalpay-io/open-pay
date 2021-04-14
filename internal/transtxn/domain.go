package transtxn

import (
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type TransferTransaction struct {
	ID               uint64
	TransferID       uint64
	SourceID         string
	CustomerID       uint64
	DestinationID    uint64
	WalletPID        uint64
	GatewayRequestID uint64
	Type             Type
	Amount           int64
	Currency         string
	Status           tcc.STATUS
	ErrorCode        string
	ErrorMsg         string
	Metadata         map[string]interface{}
	CreatedAt        int64
	UpdatedAt        int64
}

type Type string

const (
	WALLET_PAY          Type = "WALLET_PAY"
	WALLET_PAY_EXTERNAL Type = "WALLET_PAY_EXTERNAL"
	CC_DIRECT           Type = "CC_DIRECT"
)

var types = [...]string{
	"WALLET_PAY",
	"WALLET_PAY_EXTERNAL",
	"CC_DIRECT",
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
