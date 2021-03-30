package transtxn

import (
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type CCTransferTransaction struct {
	TransferTransaction
	Observer tcc.Observer
	repo     Repository
	// FIXME: add gateway service here
	// Gateway
}

func (c *CCTransferTransaction) Try() error {
	return nil
}

func (c *CCTransferTransaction) Commit() error {
	return nil
}

func (c *CCTransferTransaction) Cancel() error {
	return nil
}

func (c *CCTransferTransaction) GetStatus() domain.STATUS {
	return c.Status
}
