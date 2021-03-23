package storage

import (
	"context"
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

func (t *topupModel) model(topup domain.Topup) {
	t.ID = topup.ID
	t.CustomerID = topup.CustomerID
	t.PaymentMethodID = topup.PaymentMethodID
	t.Amount = topup.Amount
	t.Currency = topup.Currency
	t.Status = topup.Status.String()
	t.Metadata = topup.Metadata
}

func (t *topupModel) domain() domain.Topup {
	return domain.Topup{
		ID:              t.ID,
		CustomerID:      t.CustomerID,
		PaymentMethodID: t.PaymentMethodID,
		Amount:          t.Amount,
		Currency:        t.Currency,
		Status:          domain.STATUS(t.Status),
		Metadata:        t.Metadata,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}

func (r *Repository) AddTopup(ctx context.Context, topup domain.Topup) (domain.Topup, error) {
	db := r.DB(ctx)

	var t topupModel
	t.model(topup)

	now := time.Now().Unix()
	t.CreatedAt = now
	t.UpdatedAt = now

	err := db.Create(&t).Error
	if err != nil {
		return domain.Topup{}, wrapDBError(err, "topup")
	}

	topup.CreatedAt = t.CreatedAt
	topup.UpdatedAt = t.UpdatedAt

	return topup, nil
}

func (r *Repository) UpdateTopup(ctx context.Context, topup domain.Topup) (domain.Topup, error) {
	db := r.DB(ctx)

	var t topupModel
	t.model(topup)

	t.UpdatedAt = time.Now().Unix()

	err := db.Model(&t).Updates(&t).Error
	if err != nil {
		return domain.Topup{}, wrapDBError(err, "topup")
	}

	topup.UpdatedAt = t.UpdatedAt

	return topup, nil
}

func (r *Repository) GetTopup(ctx context.Context, topupID uint64) (domain.Topup, error) {
	db := r.DB(ctx)

	var t topupModel
	err := db.Unscoped().First(&t, topupID).Error
	if err != nil {
		return domain.Topup{}, wrapDBError(err, "topup")
	}

	topup := t.domain()

	return topup, nil
}
