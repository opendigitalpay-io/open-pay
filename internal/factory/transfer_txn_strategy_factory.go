package factory

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
	"github.com/opendigitalpay-io/open-pay/internal/trans"
	"github.com/opendigitalpay-io/open-pay/internal/transtxn"
)

type transferTxnStrategyFactory struct {
	service transtxn.Service
}

type TransferTxnStrategyFactory interface {
	CreateCCDirectTransferTxnStrategy(context.Context, trans.Transfer) (tcc.Strategy, error)
	CreateBalancePayTransferTxnStrategy(context.Context, trans.Transfer) (tcc.Strategy, error)
	CreateBalanceExternalPayTransferTxnStrategy(context.Context, trans.Transfer) (tcc.Strategy, error)
}

func NewTransferTxnStrategyFactory(service transtxn.Service) TransferTxnStrategyFactory {
	return &transferTxnStrategyFactory{
		service: service,
	}
}

func (c *transferTxnStrategyFactory) CreateCCDirectTransferTxnStrategy(ctx context.Context, transfer trans.Transfer) (tcc.Strategy, error) {
	strategy, err := c.create(ctx, transfer, transtxn.CC_DIRECT)
	if err != nil {
		return &transtxn.CCTransferTransactionStrategy{}, err
	}

	return strategy, nil
}

func (c *transferTxnStrategyFactory) CreateBalancePayTransferTxnStrategy(ctx context.Context, transfer trans.Transfer) (tcc.Strategy, error) {
	strategy, err := c.create(ctx, transfer, transtxn.WALLET_PAY)
	if err != nil {
		return &transtxn.BalanceTransferTransactionStrategy{}, err
	}

	return strategy, nil
}

func (c *transferTxnStrategyFactory) CreateBalanceExternalPayTransferTxnStrategy(ctx context.Context, transfer trans.Transfer) (tcc.Strategy, error) {
	strategy, err := c.create(ctx, transfer, transtxn.WALLET_PAY_EXTERNAL)
	if err != nil {
		return &transtxn.BalanceTransferTransactionStrategy{}, err
	}

	return strategy, nil
}

func (c *transferTxnStrategyFactory) create(ctx context.Context, transfer trans.Transfer, transactionType transtxn.Type) (tcc.Strategy, error) {
	transTxn := transtxn.TransferTransaction{
		TransferID:    transfer.ID,
		SourceID:      transfer.SourceID,
		CustomerID:    transfer.CustomerID,
		DestinationID: transfer.DestinationID,
		Type:          transactionType,
		Amount:        transfer.Amount,
		Currency:      transfer.Currency,
		Status:        tcc.CREATED,
		Metadata:      transfer.Metadata,
	}

	transTxn, err := c.service.AddTransferTransaction(ctx, transTxn)
	if err != nil {
		return &transtxn.CCTransferTransactionStrategy{}, err
	}

	var transferTxnStrategy tcc.Strategy
	switch transactionType {
	case transtxn.CC_DIRECT:
		transferTxnStrategy = &transtxn.CCTransferTransactionStrategy{
			TransferTransaction: transTxn,
			Service:             c.service,
		}
	case transtxn.WALLET_PAY:
		transferTxnStrategy = &transtxn.BalanceTransferTransactionStrategy{
			TransferTransaction: transTxn,
			Service:             c.service,
		}
	case transtxn.WALLET_PAY_EXTERNAL:
		transferTxnStrategy = &transtxn.BalanceTransferTransactionStrategy{
			TransferTransaction: transTxn,
			Service:             c.service,
		}
	default:
		// TODO: raise exception
	}

	return transferTxnStrategy, nil
}
