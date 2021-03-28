package transfer_transaction

import (
	"github.com/opendigitalpay-io/open-pay/internal/domain"
)

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

type TransferTransactionDTO struct {
	TransferID    uint64
	SourceID      uint64
	DestinationID uint64
	Type          string
	Amount        int64
	Currency      string
	Metadata      map[string]interface{}
}

func Create(dto TransferTransactionDTO) (TransferTransaction, error) {
	return TransferTransaction{
		TransferID:    dto.TransferID,
		SourceID:      dto.SourceID,
		DestinationID: dto.DestinationID,
		Type:          dto.Type,
		Amount:        dto.Amount,
		Currency:      dto.Currency,
		Metadata:      dto.Metadata,
	}, nil
}
