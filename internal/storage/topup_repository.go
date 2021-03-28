package storage

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"time"
)

type topupModel struct {
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

func (t *topupModel) TableName() string {
	return "topups"
}

func (t *topupModel) model(topup domain.Topup) error {
	t.ID = topup.ID
	t.CustomerID = topup.CustomerID
	t.PaymentMethodID = topup.PaymentMethodID
	t.Amount = topup.Amount
	t.Currency = topup.Currency
	t.Status = topup.Status.String()
	metadata, err := jsoniter.Marshal(topup.Metadata)
	if err != nil {
		return err
	}

	t.Metadata = metadata

	return nil
}

func (t *topupModel) domain() (domain.Topup, error) {
	var metadata map[string]interface{}
	err := jsoniter.Unmarshal(t.Metadata, &metadata)
	if err != nil {
		return domain.Topup{}, err
	}

	return domain.Topup{
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

func (r *Repository) AddTopUp(ctx context.Context, topup domain.Topup) (domain.Topup, error) {
	db := r.DB(ctx)

	var t topupModel
	err := t.model(topup)
	if err != nil {
		return domain.Topup{}, wrapDBError(err, "topup")
	}

	now := time.Now().Unix()
	t.CreatedAt = now
	t.UpdatedAt = now

	err = db.Create(&t).Error
	if err != nil {
		return domain.Topup{}, wrapDBError(err, "topup")
	}

	topup.CreatedAt = t.CreatedAt
	topup.UpdatedAt = t.UpdatedAt

	return topup, nil
}

func (r *Repository) UpdateTopUp(ctx context.Context, topup domain.Topup) (domain.Topup, error) {
	db := r.DB(ctx)

	var t topupModel
	err := t.model(topup)
	if err != nil {
		return domain.Topup{}, wrapDBError(err, "topup")
	}

	t.UpdatedAt = time.Now().Unix()

	err = db.Model(&t).Updates(&t).Error
	if err != nil {
		return domain.Topup{}, wrapDBError(err, "topup")
	}

	topup.UpdatedAt = t.UpdatedAt

	return topup, nil
}

func (r *Repository) GetTopUp(ctx context.Context, topupID uint64) (domain.Topup, error) {
	db := r.DB(ctx)

	var t topupModel
	err := db.Unscoped().First(&t, topupID).Error
	if err != nil {
		return domain.Topup{}, wrapDBError(err, "topup")
	}

	topup, err := t.domain()
	if err != nil {
		return domain.Topup{}, wrapDBError(err, "topup")
	}

	return topup, nil
}
