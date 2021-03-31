package transtxn

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
	"github.com/opendigitalpay-io/open-pay/internal/trans"
)

type strategyFactory struct {
	repo Repository
	//GatewayService
}

type StrategyFactory interface {
	Create(context.Context, trans.TransferStrategy) (tcc.Strategy, error)
}

func NewStrategyFactory() StrategyFactory {
	return &strategyFactory{}
}

func (c *strategyFactory) Create(ctx context.Context, transfer trans.TransferStrategy) (tcc.Strategy, error) {
	transTxn := TransferTransaction{
		TransferID:    transfer.ID,
		SourceID:      transfer.SourceID,
		DestinationID: transfer.DestinationID,
		Type:          transfer.Type.String(), // FIXME: add logic to determine Type: WALLET_TOPUP, CCDIRECT, CCREFUND
		Amount:        transfer.Amount,
		Currency:      transfer.Currency,
		Status:        domain.CREATED,
		Metadata:      transfer.Metadata,
	}

	transTxn, err := c.repo.AddTransferTransaction(ctx, transTxn)
	if err != nil {
		return &CCTransferTransactionStrategy{}, err
	}

	// FIXME: add if else logic to determine what transfer txn.
	ccTransferTxn := CCTransferTransactionStrategy{
		TransferTransaction: transTxn,
		transferObserver:    &transfer,
	}

	return &ccTransferTxn, nil
}
