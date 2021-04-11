package transtxn

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
	"github.com/opendigitalpay-io/open-pay/internal/trans"
)

type strategyFactory struct {
	service Service
}

type StrategyFactory interface {
	CreateCCDirectTransferTxnStrategy(context.Context, trans.Transfer) (tcc.Strategy, error)
	CreateWalletPayTransferTxnStrategy(context.Context, trans.Transfer) (tcc.Strategy, error)
}

func NewStrategyFactory(service Service) StrategyFactory {
	return &strategyFactory{
		service: service,
	}
}

func (c *strategyFactory) CreateCCDirectTransferTxnStrategy(ctx context.Context, transfer trans.Transfer) (tcc.Strategy, error) {
	strategy, err := c.create(ctx, transfer, CC_DIRECT)
	if err != nil {
		return &CCTransferTransactionStrategy{}, err
	}

	return strategy, nil
}

func (c *strategyFactory) CreateWalletPayTransferTxnStrategy(ctx context.Context, transfer trans.Transfer) (tcc.Strategy, error) {
	strategy, err := c.create(ctx, transfer, WALLET_PAY)
	if err != nil {
		return &CCTransferTransactionStrategy{}, err
	}

	return strategy, nil
}

func (c *strategyFactory) create(ctx context.Context, transfer trans.Transfer, transactionType Type) (tcc.Strategy, error) {
	transTxn := TransferTransaction{
		TransferID:    transfer.ID,
		SourceID:      transfer.SourceID,
		CustomerID:    transfer.CustomerID,
		DestinationID: transfer.DestinationID,
		Type:          transactionType.String(),
		Amount:        transfer.Amount,
		Currency:      transfer.Currency,
		Status:        tcc.CREATED,
		Metadata:      transfer.Metadata,
	}

	transTxn, err := c.service.AddTransferTransaction(ctx, transTxn)
	if err != nil {
		return &CCTransferTransactionStrategy{}, err
	}

	var transferTxnStrategy tcc.Strategy
	switch transactionType {
	case CC_DIRECT:
		transferTxnStrategy = &CCTransferTransactionStrategy{
			TransferTransaction: transTxn,
			service:             c.service,
		}
	case WALLET_PAY:
	default:
		transferTxnStrategy = &BalanceTransferTransactionStrategy{
			TransferTransaction: transTxn,
			service:             c.service,
		}
	}

	return transferTxnStrategy, nil
}
