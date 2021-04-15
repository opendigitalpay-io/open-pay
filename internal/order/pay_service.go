package order

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/port/api"
)

type PayService interface {
	PayOrder(context.Context, uint64, api.PayOrderRequest) error
}

type payService struct {
	strategyFactory StrategyFactory
}

func NewPayService(strategyFactory StrategyFactory) PayService {
	return &payService{
		strategyFactory: strategyFactory,
	}
}

func (s *payService) PayOrder(ctx context.Context, orderID uint64, request api.PayOrderRequest) error {
	paymentSource := domain.PaymentSource{
		Type: domain.PaymentSourceType(request.PaymentSource.Type),
		ID:   request.PaymentSource.ID,
	}

	orderStrategy, err := s.strategyFactory.Create(ctx, orderID, paymentSource)
	if err != nil {
		return err
	}

	tryErr := orderStrategy.Try(ctx)
	if tryErr != nil {
		orderStrategy.Cancel(ctx)
		return tryErr
	}

	return orderStrategy.Commit(ctx)
}
