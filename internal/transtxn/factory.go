package transtxn

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
	"github.com/opendigitalpay-io/open-pay/internal/trans"
)

type factory struct {
	repo Repository
	//GatewayService
}

type Factory interface {
	Create(context.Context, trans.Transfer) (tcc.Interface, error)
}

func NewFactory() Factory {
	return &factory{}
}

func (c *factory) Create(ctx context.Context, transfer trans.Transfer) (tcc.Interface, error) {
	transTxn := TransferTransaction{
		TransferID:    transfer.ID,
		SourceID:      transfer.SourceID,
		DestinationID: transfer.DestinationID,
		Type:          transfer.Type.String(), // FIXME: add logic to determine Type: WALLET_TOPUP, CCDIRECT, CCREFUND
		Amount:        transfer.Amount,
		Currency:      transfer.Currency,
		Status:        domain.CREATED,
		Metadata:      transfer.Metadata,
	}

	transTxn, err := c.repo.AddTransferTransaction(ctx, transTxn)
	if err != nil {
		return &CCTransferTransaction{}, err
	}

	// FIXME: add if else logic to determine what transfer txn.
	ccTransferTxn := CCTransferTransaction{
		transTxn,
		&transfer,
		c.repo,
	}

	return &ccTransferTxn, nil
}
