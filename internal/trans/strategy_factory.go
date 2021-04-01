package trans

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
	"github.com/opendigitalpay-io/open-pay/internal/transtxn"
)

type StrategyFactory interface {
	CreateByOrder(context.Context, domain.Order) (tcc.Strategy, error)
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

func (f *strategyFactory) CreateByOrder(ctx context.Context, order domain.Order) (tcc.Strategy, error) {
	transfer, err := f.service.AddTransfer(ctx, order)
	if err != nil {
		return &Strategy{}, err
	}

	ccTransferTxnStrategy, err := f.transferTxnStrategyFactory.CreateCCDirectTransferTxnStrategy(ctx, transfer)
	if err != nil {
		return &Strategy{}, err
	}
	balanceTransferTxnStrategy, err := f.transferTxnStrategyFactory.CreateWalletPayTransferTxnStrategy(ctx, transfer)
	if err != nil {
		return &Strategy{}, err
	}

	var transferTxnStrategies []tcc.Strategy
	transferTxnStrategies = append(transferTxnStrategies, ccTransferTxnStrategy)
	transferTxnStrategies = append(transferTxnStrategies, balanceTransferTxnStrategy)

	strategy := Strategy{
		Transfer:              transfer,
		transferTxnStrategies: transferTxnStrategies,
		service:               f.service,
	}
	ccTransferTxnStrategy.AddObserver(&strategy)
	balanceTransferTxnStrategy.AddObserver(&strategy)

	return &strategy, nil
}
