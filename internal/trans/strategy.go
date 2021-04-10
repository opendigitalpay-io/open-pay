package trans

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type TransferStrategy struct {
	Transfer
	// FIXME: add transfer service
	orderObserver         tcc.Observer
	transferTxnStrategies []tcc.Strategy
}

func (c *TransferStrategy) Try() error {
	//resp: = c.service.callGateway(c)
	//c.Observer.OnCancelFailCallback()
	return nil
}

func (c *TransferStrategy) Commit() error {
	return nil
}

func (c *TransferStrategy) Cancel() error {
	return nil
}

func (c *TransferStrategy) GetStatus() domain.STATUS {
	return c.Status
}

func (c *TransferStrategy) OnTrySuccessCallback(ctx context.Context) {

}

func (c *TransferStrategy) OnTryFailCallback(ctx context.Context) {

}

func (c *TransferStrategy) OnCommitSuccessCallback(ctx context.Context) {

}

func (c *TransferStrategy) OnCommitFailCallback(ctx context.Context) {

}

func (c *TransferStrategy) OnCancelSuccessCallback(ctx context.Context) {

}

func (c *TransferStrategy) OnCancelFailCallback(ctx context.Context) {

}
