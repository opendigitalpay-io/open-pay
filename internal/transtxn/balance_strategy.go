package transtxn

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type BalanceTransferTransactionStrategy struct {
	TransferTransaction
	transferObserver tcc.Observer
	service          Service
}

func (c *BalanceTransferTransactionStrategy) GetStatus() tcc.STATUS {
	return c.Status
}

func (c *BalanceTransferTransactionStrategy) AddObserver(observer tcc.Observer) {
	c.transferObserver = observer
}

func (c *BalanceTransferTransactionStrategy) Try(ctx context.Context) error {
	_, err := c.service.TryWalletPay(ctx, c.TransferTransaction)
	if err != nil {
		c.transferObserver.OnTryFailCallback(ctx)
		return err
	}

	c.transferObserver.OnTrySuccessCallback(ctx)
	return nil
}

func (c *BalanceTransferTransactionStrategy) Commit(ctx context.Context) error {
	_, err := c.service.CommitWalletPay(ctx, c.TransferTransaction)
	if err != nil {
		c.transferObserver.OnCommitFailCallback(ctx)
		return err
	}

	c.transferObserver.OnCommitSuccessCallback(ctx)
	return nil
}

func (c *BalanceTransferTransactionStrategy) Cancel(ctx context.Context) error {
	_, err := c.service.CancelWalletPay(ctx, c.TransferTransaction)
	if err != nil {
		c.transferObserver.OnCancelFailCallback(ctx)
		return err
	}

	c.transferObserver.OnCancelSuccessCallback(ctx)
	return nil
}
