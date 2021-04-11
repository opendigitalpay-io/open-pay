package order

import "context"

type PayService interface {
	PayOrder(context.Context, uint64) error
}

type payService struct {
	strategyFactory StrategyFactory
}

func NewPayService(strategyFactory StrategyFactory) PayService {
	return &payService{
		strategyFactory: strategyFactory,
	}
}

// FIXME: add pay-order request dto and response
func (s *payService) PayOrder(ctx context.Context, orderID uint64) error {
	strategy, err := s.strategyFactory.Create(ctx, orderID)
	if err != nil {
		return err
	}

	tryErr := strategy.Try(ctx)
	if tryErr != nil {
		strategy.Cancel(ctx)
		return tryErr
	}

	return strategy.Commit(ctx)
}
