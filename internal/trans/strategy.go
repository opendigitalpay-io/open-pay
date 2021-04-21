package trans

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type Strategy struct {
	Transfer
	OrderObserver         tcc.Observer
	TransferTxnStrategies []tcc.Strategy
	Service               Service
}

func (s *Strategy) GetStatus() tcc.STATUS {
	return s.Status
}

func (s *Strategy) AddObserver(observer tcc.Observer) {
	s.OrderObserver = observer
}

func (s *Strategy) Try(ctx context.Context) error {
	s.Status = tcc.TRY_STARTED
	s.Service.UpdateTransfer(ctx, s.Transfer)

	for _, transferTxnStrategy := range s.TransferTxnStrategies {
		err := transferTxnStrategy.Try(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Strategy) Commit(ctx context.Context) error {
	s.Status = tcc.COMMIT_STARTED
	s.Service.UpdateTransfer(ctx, s.Transfer)

	for _, transferTxnStrategy := range s.TransferTxnStrategies {
		err := transferTxnStrategy.Commit(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Strategy) Cancel(ctx context.Context) error {
	s.Status = tcc.CANCEL_STARTED
	s.Service.UpdateTransfer(ctx, s.Transfer)

	for _, transferTxnStrategy := range s.TransferTxnStrategies {
		err := transferTxnStrategy.Cancel(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Strategy) OnTrySuccessCallback(ctx context.Context) {
	if s.allTransferTxnStrategiesMatch(tcc.TRY_SUCCEEDED) {
		s.Status = tcc.TRY_SUCCEEDED
		s.Service.UpdateTransfer(ctx, s.Transfer)
		s.OrderObserver.OnTrySuccessCallback(ctx)
	}
}

func (s *Strategy) OnTryFailCallback(ctx context.Context) {
	s.Status = tcc.TRY_FAILED
	s.Service.UpdateTransfer(ctx, s.Transfer)
	s.OrderObserver.OnTryFailCallback(ctx)
}

func (s *Strategy) OnCommitSuccessCallback(ctx context.Context) {
	if s.allTransferTxnStrategiesMatch(tcc.COMMIT_SUCCEEDED) {
		s.Status = tcc.COMMIT_SUCCEEDED
		s.Service.UpdateTransfer(ctx, s.Transfer)
		s.OrderObserver.OnCommitSuccessCallback(ctx)
	}
}

func (s *Strategy) OnCommitFailCallback(ctx context.Context) {
	s.Status = tcc.COMMIT_FAILED
	s.Service.UpdateTransfer(ctx, s.Transfer)
	s.OrderObserver.OnCommitFailCallback(ctx)
}

func (s *Strategy) OnCancelSuccessCallback(ctx context.Context) {
	if s.allTransferTxnStrategiesMatch(tcc.CANCEL_SUCCEEDED) {
		s.Status = tcc.CANCEL_SUCCEEDED
		s.Service.UpdateTransfer(ctx, s.Transfer)
		s.OrderObserver.OnCancelSuccessCallback(ctx)
	}
}

func (s *Strategy) OnCancelFailCallback(ctx context.Context) {
	s.Status = tcc.CANCEL_FAILED
	s.Service.UpdateTransfer(ctx, s.Transfer)
	s.OrderObserver.OnCancelFailCallback(ctx)
}

func (s *Strategy) allTransferTxnStrategiesMatch(status tcc.STATUS) bool {
	for _, transferTxnStrategy := range s.TransferTxnStrategies {
		if transferTxnStrategy.GetStatus() != status {
			return false
		}
	}
	return true
}
