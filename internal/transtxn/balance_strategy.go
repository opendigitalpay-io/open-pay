package transtxn

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type BalanceTransferTransactionStrategy struct {
	TransferTransaction
	TransferObserver tcc.Observer
	Service          Service
}

func (c *BalanceTransferTransactionStrategy) GetStatus() tcc.STATUS {
	return c.Status
}

func (c *BalanceTransferTransactionStrategy) AddObserver(observer tcc.Observer) {
	c.TransferObserver = observer
}

func (c *BalanceTransferTransactionStrategy) Try(ctx context.Context) error {
	c.Status = tcc.TRY_STARTED
	c.Service.UpdateTransferTransaction(ctx, c.TransferTransaction)

	transferTransaction, err := c.Service.TryWalletPay(ctx, c.TransferTransaction)
	if err != nil {
		c.Status = tcc.TRY_FAILED
		c.Service.UpdateTransferTransaction(ctx, c.TransferTransaction)

		c.TransferObserver.OnTryFailCallback(ctx)
		return err
	}

	// FIXME: move back to line 27
	c.TransferTransaction = transferTransaction
	c.Status = tcc.TRY_SUCCEEDED
	c.Service.UpdateTransferTransaction(ctx, c.TransferTransaction)

	c.TransferObserver.OnTrySuccessCallback(ctx)
	return nil
}

func (c *BalanceTransferTransactionStrategy) Commit(ctx context.Context) error {
	c.Status = tcc.COMMIT_STARTED
	c.Service.UpdateTransferTransaction(ctx, c.TransferTransaction)

	transferTransaction, err := c.Service.CommitWalletPay(ctx, c.TransferTransaction)
	if err != nil {
		c.Status = tcc.COMMIT_FAILED
		c.Service.UpdateTransferTransaction(ctx, c.TransferTransaction)

		c.TransferObserver.OnCommitFailCallback(ctx)
		return err
	}

	c.TransferTransaction = transferTransaction
	c.Status = tcc.COMMIT_SUCCEEDED
	c.Service.UpdateTransferTransaction(ctx, c.TransferTransaction)

	c.TransferObserver.OnCommitSuccessCallback(ctx)
	return nil
}

func (c *BalanceTransferTransactionStrategy) Cancel(ctx context.Context) error {
	c.Status = tcc.CANCEL_STARTED
	c.Service.UpdateTransferTransaction(ctx, c.TransferTransaction)

	transferTransaction, err := c.Service.CancelWalletPay(ctx, c.TransferTransaction)
	if err != nil {
		c.Status = tcc.CANCEL_FAILED
		c.Service.UpdateTransferTransaction(ctx, c.TransferTransaction)

		c.TransferObserver.OnCancelFailCallback(ctx)
		return err
	}

	c.TransferTransaction = transferTransaction
	c.Status = tcc.CANCEL_SUCCEEDED
	c.Service.UpdateTransferTransaction(ctx, c.TransferTransaction)

	c.TransferObserver.OnCancelSuccessCallback(ctx)
	return nil
}
