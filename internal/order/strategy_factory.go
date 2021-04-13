package order

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
	"github.com/opendigitalpay-io/open-pay/internal/trans"
)

type StrategyFactory interface {
	Create(context.Context, uint64) (tcc.Strategy, error)
}

type strategyFactory struct {
	service                 Service
	transferStrategyFactory trans.StrategyFactory
}

func NewStrategyFactory(service Service, transferStrategyFactory trans.StrategyFactory) StrategyFactory {
	return &strategyFactory{
		service:                 service,
		transferStrategyFactory: transferStrategyFactory,
	}
}

func (f *strategyFactory) Create(ctx context.Context, orderID uint64) (tcc.Strategy, error) {
	order, err := f.service.GetOrder(ctx, orderID)
	if err != nil {
		return &Strategy{}, err
	}

	transferStrategy, err := f.transferStrategyFactory.CreateByOrder(ctx, order)
	if err != nil {
		return &Strategy{}, err
	}

	strategy := Strategy{
		Order:            order,
		transferStrategy: transferStrategy,
		service:          f.service,
	}
	transferStrategy.AddObserver(&strategy)

	return &strategy, nil
}
