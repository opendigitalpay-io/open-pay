package transtxn

import (
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type CCTransferTransactionStrategy struct {
	TransferTransaction
	transferObserver tcc.Observer
	// FIXME: add TransferTxnService
}

func (c *CCTransferTransactionStrategy) Try() error {
	//resp: = c.service.callGateway(c)
	//c.Observer.OnCancelFailCallback()
	return nil
}

func (c *CCTransferTransactionStrategy) Commit() error {
	return nil
}

func (c *CCTransferTransactionStrategy) Cancel() error {
	return nil
}

func (c *CCTransferTransactionStrategy) GetStatus() domain.STATUS {
	return c.Status
}
