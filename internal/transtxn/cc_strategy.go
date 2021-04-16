package transtxn

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type CCTransferTransactionStrategy struct {
	TransferTransaction
	TransferObserver tcc.Observer
	Service          Service
}

func (c *CCTransferTransactionStrategy) GetStatus() tcc.STATUS {
	return c.Status
}

func (c *CCTransferTransactionStrategy) AddObserver(observer tcc.Observer) {
	c.TransferObserver = observer
}

func (c *CCTransferTransactionStrategy) Try(ctx context.Context) error {
	c.Status = tcc.TRY_STARTED
	c.Service.UpdateTransferTransaction(ctx, c.TransferTransaction)

	transferTransaction, err := c.Service.CCAuth(ctx, c.TransferTransaction)
	if err != nil {
		c.Status = tcc.TRY_FAILED
		c.Service.UpdateTransferTransaction(ctx, c.TransferTransaction)

		c.TransferObserver.OnTrySuccessCallback(ctx)
		return err
	}

	c.TransferTransaction = transferTransaction
	c.Status = tcc.TRY_SUCCEEDED
	c.Service.UpdateTransferTransaction(ctx, c.TransferTransaction)

	c.TransferObserver.OnTrySuccessCallback(ctx)
	return nil
}

func (c *CCTransferTransactionStrategy) Commit(ctx context.Context) error {
	c.Status = tcc.COMMIT_STARTED
	c.Service.UpdateTransferTransaction(ctx, c.TransferTransaction)

	transferTransaction, err := c.Service.CCCapture(ctx, c.TransferTransaction)
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

func (c *CCTransferTransactionStrategy) Cancel(ctx context.Context) error {
	c.Status = tcc.CANCEL_STARTED
	c.Service.UpdateTransferTransaction(ctx, c.TransferTransaction)

	// TODO: CCRefund
	//transferTransaction, err := c.TransferTransaction, nil
	//c.TransferTransaction = transferTransaction
	//if err != nil {
	//	c.Status = tcc.CANCEL_FAILED
	//	c.Service.UpdateTransferTransaction(ctx, c.TransferTransaction)
	//
	//	c.TransferObserver.OnCancelFailCallback(ctx)
	//	return err
	//}

	c.Status = tcc.CANCEL_SUCCEEDED
	c.Service.UpdateTransferTransaction(ctx, c.TransferTransaction)

	c.TransferObserver.OnCancelSuccessCallback(ctx)
	return nil
}
