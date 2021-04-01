package trans

import (
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

func (c *TransferStrategy) OnTrySuccessCallback() {

}

func (c *TransferStrategy) OnTryFailCallback() {

}

func (c *TransferStrategy) OnCommitSuccessCallback() {

}

func (c *TransferStrategy) OnCommitFailCallback() {

}

func (c *TransferStrategy) OnCancelSuccessCallback() {

}

func (c *TransferStrategy) OnCancelFailCallback() {

}
