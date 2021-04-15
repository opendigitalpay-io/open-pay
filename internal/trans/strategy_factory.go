package trans

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
	"github.com/opendigitalpay-io/open-pay/internal/transtxn"
)

type StrategyFactory interface {
	Create(context.Context, domain.Order, domain.PaymentSource) (tcc.Strategy, error)
}

type strategyFactory struct {
	service                    Service
	transferTxnStrategyFactory transtxn.StrategyFactory
}

func NewStrategyFactory(service Service, transferTxnStrategyFactory transtxn.StrategyFactory) StrategyFactory {
	return &strategyFactory{
		service:                    service,
		transferTxnStrategyFactory: transferTxnStrategyFactory,
	}
}

func (f *strategyFactory) Create(ctx context.Context, order domain.Order, paymentSource domain.PaymentSource) (tcc.Strategy, error) {
	transfer, err := f.service.AddTransfer(ctx, order, paymentSource)
	if err != nil {
		return &Strategy{}, err
	}

	var transferTxnStrategies []tcc.Strategy
	switch paymentSource.Type {
	case domain.TOKEN:
		ccTransferTxnStrategy, err := f.transferTxnStrategyFactory.CreateCCDirectTransferTxnStrategy(ctx, transfer)
		if err != nil {
			return &Strategy{}, err
		}
		balanceTransferTxnStrategy, err := f.transferTxnStrategyFactory.CreateBalanceExternalPayTransferTxnStrategy(ctx, transfer)
		if err != nil {
			return &Strategy{}, err
		}
		transferTxnStrategies = append(transferTxnStrategies, ccTransferTxnStrategy)
		transferTxnStrategies = append(transferTxnStrategies, balanceTransferTxnStrategy)
	case domain.BALANCE_ACCOUNT:
		balanceTransferTxnStrategy, err := f.transferTxnStrategyFactory.CreateBalancePayTransferTxnStrategy(ctx, transfer)
		if err != nil {
			return &Strategy{}, err
		}
		transferTxnStrategies = append(transferTxnStrategies, balanceTransferTxnStrategy)
	// FIXME: add INTERACT & PAYMENT_METHOD
	}

	strategy := Strategy{
		Transfer:              transfer,
		transferTxnStrategies: transferTxnStrategies,
		service:               f.service,
	}
	for _, transferTxnStrategy := range transferTxnStrategies {
		transferTxnStrategy.AddObserver(&strategy)
	}

	return &strategy, nil
}
