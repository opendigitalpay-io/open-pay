package order

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type Strategy struct {
	domain.Order
	transferStrategy tcc.Strategy
	service          Service
}

func (s *Strategy) GetStatus() tcc.STATUS {
	return s.Status
}

func (s *Strategy) AddObserver(tcc.Observer) {
}

func (s *Strategy) Try(ctx context.Context) error {
	s.Status = tcc.TRY_STARTED
	_, err := s.service.UpdateOrder(ctx, s.Order)
	if err != nil {
		return err
	}

	return s.transferStrategy.Try(ctx)
}

func (s *Strategy) Commit(ctx context.Context) error {
	s.Status = tcc.COMMIT_STARTED
	_, err := s.service.UpdateOrder(ctx, s.Order)
	if err != nil {
		return err
	}

	return s.transferStrategy.Commit(ctx)
}

func (s *Strategy) Cancel(ctx context.Context) error {
	s.Status = tcc.CANCEL_STARTED
	_, err := s.service.UpdateOrder(ctx, s.Order)
	if err != nil {
		return err
	}

	return s.transferStrategy.Cancel(ctx)
}

func (s *Strategy) OnTrySuccessCallback(ctx context.Context) {
	s.Status = tcc.TRY_SUCCEEDED
	s.service.UpdateOrder(ctx, s.Order)
}

func (s *Strategy) OnTryFailCallback(ctx context.Context) {
	s.Status = tcc.TRY_FAILED
	s.service.UpdateOrder(ctx, s.Order)
}

func (s *Strategy) OnCommitSuccessCallback(ctx context.Context) {
	s.Status = tcc.COMMIT_SUCCEEDED
	s.service.UpdateOrder(ctx, s.Order)
}

func (s *Strategy) OnCommitFailCallback(ctx context.Context) {
	s.Status = tcc.COMMIT_FAILED
	s.service.UpdateOrder(ctx, s.Order)
}

func (s *Strategy) OnCancelSuccessCallback(ctx context.Context) {
	s.Status = tcc.CANCEL_SUCCEEDED
	s.service.UpdateOrder(ctx, s.Order)
}

func (s *Strategy) OnCancelFailCallback(ctx context.Context) {
	s.Status = tcc.CANCEL_FAILED
	s.service.UpdateOrder(ctx, s.Order)
}
