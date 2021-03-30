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
	Type             string
	Amount           int64
	Currency         string
	Status           domain.STATUS
	ErrorCode        string
	ErrorMsg         string
	Metadata         map[string]interface{}
	CreatedAt        int64
	UpdatedAt        int64
}
