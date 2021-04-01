package transtxn

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
)

type Repository interface {
	AddTransferTransaction(context.Context, TransferTransaction) (TransferTransaction, error)
	UpdateTransferTransaction(context.Context, TransferTransaction) (TransferTransaction, error)
}

type TransferTransaction struct {
	ID               uint64
	TransferID       uint64
	SourceID         uint64
	DestinationID    uint64
	WalletPID        uint64
	GatewayRequestID uint64
	Type             Type
	Amount           int64
	Currency         string
	Status           domain.STATUS
	ErrorCode        string
	ErrorMsg         string
	Metadata         map[string]interface{}
	CreatedAt        int64
	UpdatedAt        int64
}

type Type string

const (
	WALLET_PAY Type = "WALLET_PAY"
	CC_DIRECT  Type = "CC_DIRECT"
)

var types = [...]string{
	"WALLET_PAY",
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
