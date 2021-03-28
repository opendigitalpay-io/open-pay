package storage

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"time"
)

type topUpModel struct {
	ID              uint64 `gorm:"primary_key"`
	CustomerID      uint64
	PaymentMethodID uint64
	Amount          int64
	Currency        string
	Status          string
	Metadata        []byte
	CreatedAt       int64
	UpdatedAt       int64
}

func (t *topUpModel) TableName() string {
	return "topups"
}

func (t *topUpModel) model(topUp domain.TopUp) error {
	t.ID = topUp.ID
	t.CustomerID = topUp.CustomerID
	t.PaymentMethodID = topUp.PaymentMethodID
	t.Amount = topUp.Amount
	t.Currency = topUp.Currency
	t.Status = topUp.Status.String()
	metadata, err := jsoniter.Marshal(topUp.Metadata)
	if err != nil {
		return err
	}

	t.Metadata = metadata

	return nil
}

func (t *topUpModel) domain() (domain.TopUp, error) {
	var metadata map[string]interface{}
	err := jsoniter.Unmarshal(t.Metadata, &metadata)
	if err != nil {
		return domain.TopUp{}, err
	}

	return domain.TopUp{
		ID:              t.ID,
		CustomerID:      t.CustomerID,
		PaymentMethodID: t.PaymentMethodID,
		Amount:          t.Amount,
		Currency:        t.Currency,
		Status:          domain.STATUS(t.Status),
		Metadata:        metadata,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}, nil
}

func (r *Repository) AddTopUp(ctx context.Context, topUp domain.TopUp) (domain.TopUp, error) {
	db := r.DB(ctx)

	var t topUpModel
	err := t.model(topUp)
	if err != nil {
		return domain.TopUp{}, wrapDBError(err, "topUp")
	}

	now := time.Now().Unix()
	t.CreatedAt = now
	t.UpdatedAt = now

	err = db.Create(&t).Error
	if err != nil {
		return domain.TopUp{}, wrapDBError(err, "topUp")
	}

	topUp.CreatedAt = t.CreatedAt
	topUp.UpdatedAt = t.UpdatedAt

	return topUp, nil
}

func (r *Repository) UpdateTopUp(ctx context.Context, topUp domain.TopUp) (domain.TopUp, error) {
	db := r.DB(ctx)

	var t topUpModel
	err := t.model(topUp)
	if err != nil {
		return domain.TopUp{}, wrapDBError(err, "topUp")
	}

	t.UpdatedAt = time.Now().Unix()

	err = db.Model(&t).Updates(&t).Error
	if err != nil {
		return domain.TopUp{}, wrapDBError(err, "topUp")
	}

	topUp.UpdatedAt = t.UpdatedAt

	return topUp, nil
}

func (r *Repository) GetTopUp(ctx context.Context, topUpID uint64) (domain.TopUp, error) {
	db := r.DB(ctx)

	var t topUpModel
	err := db.Unscoped().First(&t, topUpID).Error
	if err != nil {
		return domain.TopUp{}, wrapDBError(err, "topUp")
	}

	topUp, err := t.domain()
	if err != nil {
		return domain.TopUp{}, wrapDBError(err, "topUp")
	}

	return topUp, nil
}
