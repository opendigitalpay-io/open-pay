package factory

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	orderpkg "github.com/opendigitalpay-io/open-pay/internal/order"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type OrderStrategyFactory interface {
	Create(context.Context, uint64, domain.PaymentSource) (tcc.Strategy, error)
}

type orderStrategyFactory struct {
	service                 orderpkg.Service
	transferStrategyFactory TransferStrategyFactory
}

func NewOrderStrategyFactory(service orderpkg.Service, transferStrategyFactory TransferStrategyFactory) OrderStrategyFactory {
	return &orderStrategyFactory{
		service:                 service,
		transferStrategyFactory: transferStrategyFactory,
	}
}

func (f *orderStrategyFactory) Create(ctx context.Context, orderID uint64, paymentSource domain.PaymentSource) (tcc.Strategy, error) {
	order, err := f.service.GetOrder(ctx, orderID)
	if err != nil {
		return &orderpkg.Strategy{}, err
	}

	transferStrategy, err := f.transferStrategyFactory.Create(ctx, order, paymentSource)
	if err != nil {
		return &orderpkg.Strategy{}, err
	}

	strategy := orderpkg.Strategy{
		Order:            order,
		TransferStrategy: transferStrategy,
		Service:          f.service,
	}
	transferStrategy.AddObserver(&strategy)

	return &strategy, nil
}
