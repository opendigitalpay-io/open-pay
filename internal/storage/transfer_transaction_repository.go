package storage

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/transfer_transaction"
	"time"
)

type transferTransactionModel struct {
	ID               uint64 `gorm:"primary_key"`
	TransferID       uint64
	SourceID         uint64
	DestinationID    uint64
	WalletPID        uint64
	GatewayRequestID uint64
	Type             string
	Amount           int64
	Currency         string
	Status           string
	ErrorCode        string
	ErrorMsg         string
	Metadata         []byte
	CreatedAt        int64
	UpdatedAt        int64
}

func (t *transferTransactionModel) TableName() string {
	return "transfer_transactions"
}

func (t *transferTransactionModel) model(txn transfer_transaction.TransferTransaction) error {
	t.ID = txn.ID
	t.TransferID = txn.TransferID
	t.SourceID = txn.SourceID
	t.DestinationID = txn.DestinationID
	t.WalletPID = txn.WalletPID
	t.GatewayRequestID = txn.GatewayRequestID
	t.Type = txn.Type
	t.Amount = txn.Amount
	t.Currency = txn.Currency
	t.Status = txn.Status.String()
	t.ErrorCode = txn.ErrorCode
	t.ErrorMsg = txn.ErrorMsg
	metadata, err := jsoniter.Marshal(txn.Metadata)
	if err != nil {
		return err
	}

	t.Metadata = metadata
	return nil
}

func (t *transferTransactionModel) domain() (transfer_transaction.TransferTransaction, error) {
	var metadata map[string]interface{}
	err := jsoniter.Unmarshal(t.Metadata, &metadata)
	if err != nil {
		return transfer_transaction.TransferTransaction{}, err
	}

	return transfer_transaction.TransferTransaction{
		ID:               t.ID,
		TransferID:       t.TransferID,
		SourceID:         t.SourceID,
		DestinationID:    t.DestinationID,
		WalletPID:        t.WalletPID,
		GatewayRequestID: t.GatewayRequestID,
		Type:             t.Type,
		Amount:           t.Amount,
		Currency:         t.Currency,
		Status:           domain.STATUS(t.Status),
		ErrorCode:        t.ErrorCode,
		ErrorMsg:         t.ErrorMsg,
		Metadata:         metadata,
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
	}, nil
}

func (r *Repository) AddTransferTransaction(ctx context.Context, transaction transfer_transaction.TransferTransaction) (transfer_transaction.TransferTransaction, error) {
	db := r.DB(ctx)

	id, err := r.uidGenerator.NextID()
	if err != nil {
		return transfer_transaction.TransferTransaction{}, wrapDBError(err, "transfer_transaction")
	}
	transaction.ID = id

	var t transferTransactionModel
	err = t.model(transaction)
	if err != nil {
		return transfer_transaction.TransferTransaction{}, wrapDBError(err, "transfer_transaction")
	}

	now := time.Now().Unix()
	t.CreatedAt = now
	t.UpdatedAt = now

	err = db.Create(&t).Error
	if err != nil {
		return transfer_transaction.TransferTransaction{}, wrapDBError(err, "transfer_transaction")
	}

	transaction.CreatedAt = t.CreatedAt
	transaction.UpdatedAt = t.UpdatedAt

	return transaction, nil
}

func (r *Repository) UpdateTransferTransaction(ctx context.Context, transaction transfer_transaction.TransferTransaction) (transfer_transaction.TransferTransaction, error) {
	db := r.DB(ctx)

	var t transferTransactionModel
	err := t.model(transaction)
	if err != nil {
		return transfer_transaction.TransferTransaction{}, wrapDBError(err, "transfer_transaction")
	}

	t.UpdatedAt = time.Now().Unix()

	err = db.Model(&t).Updates(&t).Error
	if err != nil {
		return transfer_transaction.TransferTransaction{}, wrapDBError(err, "transfer_transaction")
	}

	transaction.UpdatedAt = t.UpdatedAt

	return transaction, nil
}
