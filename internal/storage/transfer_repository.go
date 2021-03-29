package storage

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/trans"
	"time"
)

type transferModel struct {
	ID        uint64 `gorm:"primary_key"`
	OrderID   uint64
	Type      string
	Status    string
	Amount    int64
	Currency  string
	Metadata  []byte
	CreatedAt int64
	UpdatedAt int64
}

func (t *transferModel) TableName() string {
	return "transfers"
}

func (t *transferModel) model(transfer trans.Transfer) error {
	t.ID = transfer.ID
	t.OrderID = transfer.OrderID
	t.Type = transfer.Type.String()
	t.Status = transfer.Status.String()
	t.Amount = transfer.Amount
	t.Currency = transfer.Currency

	metadata, err := jsoniter.Marshal(transfer.Metadata)
	if err != nil {
		return err
	}
	t.Metadata = metadata

	return nil
}

func (t *transferModel) domain() (trans.Transfer, error) {
	var metadata map[string]interface{}
	err := jsoniter.Unmarshal(t.Metadata, &metadata)
	if err != nil {
		return trans.Transfer{}, err
	}

	return trans.Transfer{
		ID:        t.ID,
		OrderID:   t.OrderID,
		Type:      trans.Type(t.Type),
		Amount:    t.Amount,
		Currency:  t.Currency,
		Status:    domain.STATUS(t.Status),
		Metadata:  metadata,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}, nil
}

func (r *Repository) AddTransfer(ctx context.Context, transfer trans.Transfer) (trans.Transfer, error) {
	db := r.DB(ctx)

	transferID, err := r.uidGenerator.NextID()
	if err != nil {
		return trans.Transfer{}, wrapDBError(err, "transfer")
	}
	transfer.ID = transferID

	var t transferModel
	err = t.model(transfer)
	if err != nil {
		return trans.Transfer{}, wrapDBError(err, "transfer")
	}

	now := time.Now().Unix()
	t.CreatedAt = now
	t.UpdatedAt = now

	err = db.Create(&t).Error
	if err != nil {
		return trans.Transfer{}, wrapDBError(err, "transfer")
	}

	transfer.CreatedAt = t.CreatedAt
	transfer.UpdatedAt = t.UpdatedAt

	return transfer, nil
}

func (r *Repository) UpdateTransfer(ctx context.Context, transfer trans.Transfer) (trans.Transfer, error) {
	db := r.DB(ctx)

	var t transferModel
	err := t.model(transfer)
	if err != nil {
		return trans.Transfer{}, wrapDBError(err, "transfer")
	}

	t.UpdatedAt = time.Now().Unix()

	err = db.Model(&t).Updates(&t).Error
	if err != nil {
		return trans.Transfer{}, wrapDBError(err, "transfer")
	}

	transfer.UpdatedAt = t.UpdatedAt

	return transfer, nil
}
