package storage

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
	"time"
)
import jsoniter "github.com/json-iterator/go"

type orderModel struct {
	ID            uint64 `gorm:"primary_key"`
	CustomerID    uint64
	MerchantID    uint64
	Amount        int64
	Currency      string
	CustomerEmail string
	ReferenceID   string
	Status        string
	Mode          string
	Metadata      []byte
	CreatedAt     int64
	UpdatedAt     int64
}

func (o *orderModel) TableName() string {
	return "orders"
}

func (o *orderModel) model(order domain.Order) error {
	o.ID = order.ID
	o.CustomerID = order.CustomerID
	o.MerchantID = order.MerchantID
	o.Amount = order.Amount
	o.Currency = order.Currency
	o.ReferenceID = order.ReferenceID
	o.Status = order.Status.String()
	o.Mode = order.Mode.String()

	metadata, err := jsoniter.Marshal(order.Metadata)
	if err != nil {
		return err
	}
	o.Metadata = metadata

	return nil
}

func (o *orderModel) domain() (domain.Order, error) {
	var metadata map[string]interface{}
	err := jsoniter.Unmarshal(o.Metadata, &metadata)
	if err != nil {
		return domain.Order{}, err
	}

	return domain.Order{
		ID:            o.ID,
		CustomerID:    o.CustomerID,
		MerchantID:    o.MerchantID,
		Amount:        o.Amount,
		Currency:      o.Currency,
		ReferenceID:   o.ReferenceID,
		CustomerEmail: o.CustomerEmail,
		Status:        tcc.STATUS(o.Status),
		Mode:          domain.OrderMode(o.Mode),
		Metadata:      metadata,
		CreatedAt:     o.CreatedAt,
		UpdatedAt:     o.UpdatedAt,
	}, nil
}

func (r *Repository) AddOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	db := r.DB(ctx)

	var o orderModel
	err := o.model(order)
	if err != nil {
		return domain.Order{}, wrapDBError(err, "order")
	}

	now := time.Now().Unix()
	o.CreatedAt = now
	o.UpdatedAt = now

	err = db.Create(&o).Error
	if err != nil {
		return domain.Order{}, wrapDBError(err, "order")
	}

	order.CreatedAt = o.CreatedAt
	order.UpdatedAt = o.UpdatedAt

	return order, nil
}

func (r *Repository) GetOrder(ctx context.Context, orderID uint64) (domain.Order, error) {
	db := r.DB(ctx)

	var o orderModel
	err := db.Unscoped().First(&o, orderID).Error
	if err != nil {
		return domain.Order{}, wrapDBError(err, "order")
	}

	order, err := o.domain()
	if err != nil {
		return domain.Order{}, wrapDBError(err, "order")
	}

	return order, nil
}

func (r *Repository) UpdateOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	db := r.DB(ctx)

	var o orderModel
	err := o.model(order)
	if err != nil {
		return domain.Order{}, wrapDBError(err, "order")
	}

	o.UpdatedAt = time.Now().Unix()

	err = db.Model(&o).Updates(&o).Error
	if err != nil {
		return domain.Order{}, wrapDBError(err, "order")
	}

	order.UpdatedAt = o.UpdatedAt

	return order, nil
}
