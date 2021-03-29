package storage

import (
	"context"
	jsoniter "github.com/json-iterator/go"
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

func (r *refundModel) model(refund domain.Refund) error {
	r.ID = refund.ID
	r.OrderID = refund.OrderID
	r.Amount = refund.Amount
	r.Status = refund.Status.String()
	r.RefundCount = refund.RefundCount
	metadata, err := jsoniter.Marshal(refund.Metadata)
	if err != nil {
		return err
	}
	r.Metadata = metadata

	return nil
}

func (r *refundModel) domain() (domain.Refund, error) {
	var metadata map[string]interface{}
	err := jsoniter.Unmarshal(r.Metadata, &metadata)
	if err != nil {
		return domain.Refund{}, err
	}

	return domain.Refund{
		ID:          r.ID,
		OrderID:     r.OrderID,
		Amount:      r.Amount,
		Status:      domain.STATUS(r.Status),
		RefundCount: r.RefundCount,
		Metadata:    metadata,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}, nil
}

func (r *Repository) AddRefund(ctx context.Context, refund domain.Refund) (domain.Refund, error) {
	db := r.DB(ctx)

	var rf refundModel
	err := rf.model(refund)
	if err != nil {
		return domain.Refund{}, wrapDBError(err, "refund")
	}

	now := time.Now().Unix()
	rf.CreatedAt = now
	rf.UpdatedAt = now

	err = db.Create(&rf).Error
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
	err := rf.model(refund)
	if err != nil {
		return domain.Refund{}, wrapDBError(err, "refund")
	}

	rf.UpdatedAt = time.Now().Unix()

	err = db.Model(&rf).Updates(&rf).Error
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

	refund, err := rf.domain()
	if err != nil {
		return domain.Refund{}, wrapDBError(err, "refund")
	}

	return refund, nil
}
