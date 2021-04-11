package transtxn

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/external/balance"
	"github.com/opendigitalpay-io/open-pay/internal/common/uid"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type Service interface {
	AddTransferTransaction(context.Context, TransferTransaction) (TransferTransaction, error)
	UpdateTransferTransaction(context.Context, TransferTransaction) (TransferTransaction, error)

	TryWalletPay(context.Context, TransferTransaction) (TransferTransaction, error)
	CommitWalletPay(context.Context, TransferTransaction) (TransferTransaction, error)
	CancelWalletPay(context.Context, TransferTransaction) (TransferTransaction, error)
}

type Repository interface {
	AddTransferTransaction(context.Context, TransferTransaction) (TransferTransaction, error)
	UpdateTransferTransaction(context.Context, TransferTransaction) (TransferTransaction, error)
}

type service struct {
	repo           Repository
	balanceAdapter balance.Adapter
	uidGenerator   uid.Generator
	// TODO: Add gateway service when it is available
}

func NewService(repo Repository, balanceAdapter balance.Adapter, uidGenerator uid.Generator) Service {
	return &service{
		repo:           repo,
		balanceAdapter: balanceAdapter,
		uidGenerator:   uidGenerator,
	}
}

func (s *service) AddTransferTransaction(ctx context.Context, transferTxn TransferTransaction) (TransferTransaction, error) {
	txn, err := s.repo.AddTransferTransaction(ctx, transferTxn)
	if err != nil {
		return TransferTransaction{}, err
	}

	return txn, nil
}

func (s *service) UpdateTransferTransaction(ctx context.Context, transferTxn TransferTransaction) (TransferTransaction, error) {
	txn, err := s.repo.UpdateTransferTransaction(ctx, transferTxn)
	if err != nil {
		return TransferTransaction{}, err
	}

	return txn, nil
}

func (s *service) TryWalletPay(ctx context.Context, transferTxn TransferTransaction) (TransferTransaction, error) {
	req := balance.TryPayRequest{
		UserID:     transferTxn.CustomerID,
		BusinessID: transferTxn.DestinationID,
		Amount:     transferTxn.Amount,
		Currency:   transferTxn.Currency,
	}

	idemKey, err := s.uidGenerator.NextID()
	if err != nil {
		return TransferTransaction{}, err
	}

	tryPayResp, err := s.balanceAdapter.TryPay(idemKey, req)
	if err != nil {
		return TransferTransaction{}, err
	}

	transferTxn.WalletPID = tryPayResp.ID
	transferTxn.Status = tcc.TRY_SUCCEEDED
	updatedTransferTxn, err := s.repo.UpdateTransferTransaction(ctx, transferTxn)
	if err != nil {
		return TransferTransaction{}, err
	}

	return updatedTransferTxn, nil
}

func (s *service) CommitWalletPay(ctx context.Context, transferTxn TransferTransaction) (TransferTransaction, error) {
	req := balance.CommitPayRequest{
		ParentID: transferTxn.WalletPID,
	}

	idemKey, err := s.uidGenerator.NextID()
	if err != nil {
		return TransferTransaction{}, err
	}

	err = s.balanceAdapter.CommitPay(idemKey, req)
	if err != nil {
		return TransferTransaction{}, err
	}

	transferTxn.Status = tcc.COMMIT_SUCCEEDED
	updatedTransferTxn, err := s.repo.UpdateTransferTransaction(ctx, transferTxn)
	if err != nil {
		return TransferTransaction{}, err
	}

	return updatedTransferTxn, nil
}

func (s *service) CancelWalletPay(ctx context.Context, transferTxn TransferTransaction) (TransferTransaction, error) {
	req := balance.CancelPayRequest{
		ParentID: transferTxn.WalletPID,
	}

	idemKey, err := s.uidGenerator.NextID()
	if err != nil {
		return TransferTransaction{}, err
	}

	err = s.balanceAdapter.CancelPay(idemKey, req)
	if err != nil {
		return TransferTransaction{}, err
	}

	transferTxn.Status = tcc.CANCEL_SUCCEEDED
	updatedTransferTxn, err := s.repo.UpdateTransferTransaction(ctx, transferTxn)
	if err != nil {
		return TransferTransaction{}, err
	}

	return updatedTransferTxn, nil
}
