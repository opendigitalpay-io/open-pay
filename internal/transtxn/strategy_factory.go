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
	CreateCCDirectTransferTxnStrategy(context.Context, trans.TransferStrategy) (tcc.Strategy, error)
	CreateWalletPayTransferTxnStrategy(context.Context, trans.TransferStrategy) (tcc.Strategy, error)
}

func NewStrategyFactory() StrategyFactory {
	return &strategyFactory{}
}

func (c *strategyFactory) CreateCCDirectTransferTxnStrategy(ctx context.Context, transfer trans.TransferStrategy) (tcc.Strategy, error) {
	strategy, err := c.create(ctx, transfer, CC_DIRECT)
	if err != nil {
		return &CCTransferTransactionStrategy{}, err
	}

	return strategy, nil
}

func (c *strategyFactory) CreateWalletPayTransferTxnStrategy(ctx context.Context, transfer trans.TransferStrategy) (tcc.Strategy, error) {
	strategy, err := c.create(ctx, transfer, WALLET_PAY)
	if err != nil {
		return &CCTransferTransactionStrategy{}, err
	}

	return strategy, nil
}

func (c *strategyFactory) create(ctx context.Context, transfer trans.TransferStrategy, transactionType Type) (tcc.Strategy, error) {
	transTxn := TransferTransaction{
		TransferID:    transfer.ID,
		SourceID:      transfer.SourceID,
		DestinationID: transfer.DestinationID,
		Type:          transactionType.String(),
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
