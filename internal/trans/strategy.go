package trans

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type Strategy struct {
	Transfer
	orderObserver         tcc.Observer
	transferTxnStrategies []tcc.Strategy
	service               Service
}

func (s *Strategy) GetStatus() tcc.STATUS {
	return s.Status
}

func (s *Strategy) AddObserver(observer tcc.Observer) {
	s.orderObserver = observer
}

func (s *Strategy) Try(ctx context.Context) error {
	s.Status = tcc.TRY_STARTED
	_, err := s.service.UpdateTransfer(ctx, s.Transfer)
	if err != nil {
		return err
	}
	for _, transferTxnStrategy := range s.transferTxnStrategies {
		err = transferTxnStrategy.Try(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Strategy) Commit(ctx context.Context) error {
	s.Status = tcc.COMMIT_STARTED
	_, err := s.service.UpdateTransfer(ctx, s.Transfer)
	if err != nil {
		return err
	}
	for _, transferTxnStrategy := range s.transferTxnStrategies {
		err = transferTxnStrategy.Commit(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Strategy) Cancel(ctx context.Context) error {
	s.Status = tcc.CANCEL_STARTED
	_, err := s.service.UpdateTransfer(ctx, s.Transfer)
	if err != nil {
		return err
	}
	for _, transferTxnStrategy := range s.transferTxnStrategies {
		err = transferTxnStrategy.Cancel(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Strategy) OnTrySuccessCallback(ctx context.Context) {
	if s.allTransferTxnStrategiesMatch(tcc.TRY_SUCCEEDED) {
		s.Status = tcc.TRY_SUCCEEDED
		s.service.UpdateTransfer(ctx, s.Transfer)
		s.orderObserver.OnTrySuccessCallback(ctx)
	}
}

func (s *Strategy) OnTryFailCallback(ctx context.Context) {
	s.Status = tcc.TRY_FAILED
	s.service.UpdateTransfer(ctx, s.Transfer)
	s.orderObserver.OnTryFailCallback(ctx)
}

func (s *Strategy) OnCommitSuccessCallback(ctx context.Context) {
	if s.allTransferTxnStrategiesMatch(tcc.COMMIT_SUCCEEDED) {
		s.Status = tcc.COMMIT_SUCCEEDED
		s.service.UpdateTransfer(ctx, s.Transfer)
		s.orderObserver.OnCommitSuccessCallback(ctx)
	}
}

func (s *Strategy) OnCommitFailCallback(ctx context.Context) {
	s.Status = tcc.COMMIT_FAILED
	s.service.UpdateTransfer(ctx, s.Transfer)
	s.orderObserver.OnCommitFailCallback(ctx)
}

func (s *Strategy) OnCancelSuccessCallback(ctx context.Context) {
	if s.allTransferTxnStrategiesMatch(tcc.CANCEL_SUCCEEDED) {
		s.Status = tcc.CANCEL_SUCCEEDED
		s.service.UpdateTransfer(ctx, s.Transfer)
		s.orderObserver.OnCancelSuccessCallback(ctx)
	}
}

func (s *Strategy) OnCancelFailCallback(ctx context.Context) {
	s.Status = tcc.CANCEL_FAILED
	s.service.UpdateTransfer(ctx, s.Transfer)
	s.orderObserver.OnCancelFailCallback(ctx)
}

func (s *Strategy) allTransferTxnStrategiesMatch(status tcc.STATUS) bool {
	for _, transferTxnStrategy := range s.transferTxnStrategies {
		if transferTxnStrategy.GetStatus() != status {
			return false
		}
	}
	return true
}
