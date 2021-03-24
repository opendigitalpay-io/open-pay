package storage

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"time"
)

type refundModel struct {
	ID          uint64 `gorm:"primary_key"`
	OrderID     uint64
	Amount      int64
	Status      string
	RefundCount int32
	Metadata    []byte
	CreatedAt   int64
	UpdatedAt   int64
}

func (t *refundModel) TableName() string {
	return "refunds"
}

func (r *refundModel) model(refund domain.Refund) {
	r.ID = refund.ID
	r.OrderID = refund.OrderID
	r.Amount = refund.Amount
	r.Status = refund.Status.String()
	r.RefundCount = refund.RefundCount
	r.Metadata = refund.Metadata
}

func (r *refundModel) domain() domain.Refund {
	return domain.Refund{
		ID:          r.ID,
		OrderID:     r.OrderID,
		Amount:      r.Amount,
		Status:      domain.STATUS(r.Status),
		RefundCount: r.RefundCount,
		Metadata:    r.Metadata,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func (r *Repository) AddRefund(ctx context.Context, refund domain.Refund) (domain.Refund, error) {
	db := r.DB(ctx)

	var rf refundModel
	rf.model(refund)

	now := time.Now().Unix()
	rf.CreatedAt = now
	rf.UpdatedAt = now

	err := db.Create(&rf).Error
	if err != nil {
		return domain.Refund{}, wrapDBError(err, "refund")
	}

	refund.CreatedAt = rf.CreatedAt
	refund.UpdatedAt = rf.UpdatedAt

	return refund, nil
}

func (r *Repository) UpdateRefund(ctx context.Context, refund domain.Refund) (domain.Refund, error) {
	db := r.DB(ctx)

	var rf refundModel
	rf.model(refund)

	rf.UpdatedAt = time.Now().Unix()

	err := db.Model(&rf).Updates(&rf).Error
	if err != nil {
		return domain.Refund{}, wrapDBError(err, "refund")
	}

	refund.UpdatedAt = rf.UpdatedAt

	return refund, nil
}

func (r *Repository) GetRefund(ctx context.Context, refundID uint64) (domain.Refund, error) {
	db := r.DB(ctx)

	var rf refundModel
	err := db.Unscoped().First(&rf, refundID).Error
	if err != nil {
		return domain.Refund{}, wrapDBError(err, "refund")
	}

	refund := rf.domain()

	return refund, nil
}
