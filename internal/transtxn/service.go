package transtxn

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/external/balance"
	"github.com/opendigitalpay-io/open-pay/internal/common/uid"
	"github.com/opendigitalpay-io/open-pay/internal/gateway"
)

type Service interface {
	AddTransferTransaction(context.Context, TransferTransaction) (TransferTransaction, error)
	UpdateTransferTransaction(context.Context, TransferTransaction) (TransferTransaction, error)

	TryWalletPay(context.Context, TransferTransaction) (TransferTransaction, error)
	CommitWalletPay(context.Context, TransferTransaction) (TransferTransaction, error)
	CancelWalletPay(context.Context, TransferTransaction) (TransferTransaction, error)

	CCAuth(context.Context, TransferTransaction) (TransferTransaction, error)
	CCCapture(context.Context, TransferTransaction) (TransferTransaction, error)
}

type Repository interface {
	AddTransferTransaction(context.Context, TransferTransaction) (TransferTransaction, error)
	UpdateTransferTransaction(context.Context, TransferTransaction) (TransferTransaction, error)
}

type service struct {
	repo           Repository
	balanceAdapter balance.Adapter
	gatewayService gateway.Service
	uidGenerator   uid.Generator
}

func NewService(repo Repository, balanceAdapter balance.Adapter, gatewayService gateway.Service, uidGenerator uid.Generator) Service {
	return &service{
		repo:           repo,
		balanceAdapter: balanceAdapter,
		gatewayService: gatewayService,
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

	var tryPayResp balance.TryPayResponse
	switch transferTxn.Type {
	case WALLET_PAY_EXTERNAL:
		tryPayResp, err = s.balanceAdapter.TryExternalPay(idemKey, req)
	case WALLET_PAY:
		tryPayResp, err = s.balanceAdapter.TryPay(idemKey, req)
	default:
	}

	if err != nil {
		return TransferTransaction{}, err
	}

	transferTxn.WalletPID = tryPayResp.ID
	return transferTxn, nil
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

	return transferTxn, nil
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

	return transferTxn, nil
}

func (s *service) CCAuth(ctx context.Context, transferTxn TransferTransaction) (TransferTransaction, error) {
	cardRequestDTO := gateway.CardRequestDTO{
		ID: transferTxn.ID,
		GatewayToken: transferTxn.SourceID,
		Amount: transferTxn.Amount,
		Currency: transferTxn.Currency,
		RequestType: gateway.AUTHORIZE,
		AutoCapture: transferTxn.Type == CC_DIRECT,
		Metadata: nil,
	}

	cardRequest, err := s.gatewayService.CreateCardRequest(ctx, cardRequestDTO)
	if err != nil {
		return TransferTransaction{}, err
	}

	_, err = s.gatewayService.Authorize(ctx, cardRequest)
	if err != nil {
		// TODO: error handling: add errorCode & errorMsg into transferTxn
		return TransferTransaction{}, err
	}

	return transferTxn, nil
}

func (s *service) CCCapture(ctx context.Context, transferTxn TransferTransaction) (TransferTransaction, error) {
	if transferTxn.Type == CC_DIRECT {
		return transferTxn, nil
	}
	cardRequestDTO := gateway.CardRequestDTO{
		ID: transferTxn.ID,
		GatewayToken: transferTxn.SourceID,
		Amount: transferTxn.Amount,
		Currency: transferTxn.Currency,
		RequestType: gateway.CAPTURE,
		AutoCapture: false,
		Metadata: nil,
	}

	cardRequest, err := s.gatewayService.CreateCardRequest(ctx, cardRequestDTO)
	if err != nil {
		return TransferTransaction{}, err
	}

	_, err = s.gatewayService.Capture(ctx, cardRequest)
	if err != nil {
		// TODO: error handling: add errorCode & errorMsg into transferTxn
		return TransferTransaction{}, err
	}

	return transferTxn, nil
}
