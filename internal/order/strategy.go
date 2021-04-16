package order

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type Strategy struct {
	domain.Order
	TransferStrategy tcc.Strategy
	Service          Service
}

func (s *Strategy) GetStatus() tcc.STATUS {
	return s.Status
}

func (s *Strategy) AddObserver(tcc.Observer) {
}

func (s *Strategy) Try(ctx context.Context) error {
	s.Status = tcc.TRY_STARTED
	s.Service.UpdateOrder(ctx, s.Order)

	return s.TransferStrategy.Try(ctx)
}

func (s *Strategy) Commit(ctx context.Context) error {
	s.Status = tcc.COMMIT_STARTED
	s.Service.UpdateOrder(ctx, s.Order)

	return s.TransferStrategy.Commit(ctx)
}

func (s *Strategy) Cancel(ctx context.Context) error {
	s.Status = tcc.CANCEL_STARTED
	s.Service.UpdateOrder(ctx, s.Order)

	return s.TransferStrategy.Cancel(ctx)
}

func (s *Strategy) OnTrySuccessCallback(ctx context.Context) {
	s.Status = tcc.TRY_SUCCEEDED
	s.Service.UpdateOrder(ctx, s.Order)
}

func (s *Strategy) OnTryFailCallback(ctx context.Context) {
	s.Status = tcc.TRY_FAILED
	s.Service.UpdateOrder(ctx, s.Order)
}

func (s *Strategy) OnCommitSuccessCallback(ctx context.Context) {
	s.Status = tcc.COMMIT_SUCCEEDED
	s.Service.UpdateOrder(ctx, s.Order)
}

func (s *Strategy) OnCommitFailCallback(ctx context.Context) {
	s.Status = tcc.COMMIT_FAILED
	s.Service.UpdateOrder(ctx, s.Order)
}

func (s *Strategy) OnCancelSuccessCallback(ctx context.Context) {
	s.Status = tcc.CANCEL_SUCCEEDED
	s.Service.UpdateOrder(ctx, s.Order)
}

func (s *Strategy) OnCancelFailCallback(ctx context.Context) {
	s.Status = tcc.CANCEL_FAILED
	s.Service.UpdateOrder(ctx, s.Order)
}
