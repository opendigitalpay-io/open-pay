package transtxn

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type CCTransferTransactionStrategy struct {
	TransferTransaction
	transferObserver tcc.Observer
	service          Service
}

func (c *CCTransferTransactionStrategy) Try(ctx context.Context) error {
	//resp: = c.service.callGateway(c)
	//c.Observer.OnCancelFailCallback()
	return nil
}

func (c *CCTransferTransactionStrategy) Commit(ctx context.Context) error {
	return nil
}

func (c *CCTransferTransactionStrategy) Cancel(ctx context.Context) error {
	return nil
}

func (c *CCTransferTransactionStrategy) GetStatus() domain.STATUS {
	return c.Status
}
