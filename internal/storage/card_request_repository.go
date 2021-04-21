package storage

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/opendigitalpay-io/open-pay/internal/gateway"
	"time"
)

type cardRequestModel struct {
	ID            uint64 `gorm:"primary_key"`
	RequestID     uint64
	ParentID      uint64
	TransferTxnID uint64
	GatewayTxnID  string
	GatewayToken  string
	Gateway       string
	Amount        int64
	Currency      string
	RequestType   string
	AutoCapture   bool
	Status        string
	Metadata      []byte
	CreatedAt     int64
	UpdatedAt     int64
}

func (c *cardRequestModel) TableName() string {
	return "card_requests"
}

func (c *cardRequestModel) model(request gateway.CardRequest) error {
	c.ID = request.ID
	c.RequestID = request.RequestID
	c.ParentID = request.ParentID
	c.TransferTxnID = request.TransferTxnID
	c.GatewayTxnID = request.GatewayTxnID
	c.GatewayToken = request.GatewayToken
	c.Gateway = request.Gateway.String()
	c.Amount = request.Amount
	c.Currency = request.Currency
	c.RequestType = request.RequestType.String()
	c.AutoCapture = request.AutoCapture
	c.Status = request.Status.String()
	metadata, err := jsoniter.Marshal(request.Metadata)
	if err != nil {
		return err
	}

	c.Metadata = metadata
	return nil
}

func (c *cardRequestModel) domain() (gateway.CardRequest, error) {
	var metadata map[string]interface{}
	err := jsoniter.Unmarshal(c.Metadata, &metadata)
	if err != nil {
		return gateway.CardRequest{}, err
	}

	return gateway.CardRequest{
		ID:            c.ID,
		RequestID:     c.RequestID,
		ParentID:      c.ParentID,
		TransferTxnID: c.TransferTxnID,
		GatewayTxnID:  c.GatewayTxnID,
		GatewayToken:  c.GatewayToken,
		Gateway:       gateway.SUPPORTED_GATEWAY(c.Gateway),
		Amount:        c.Amount,
		Currency:      c.Currency,
		RequestType:   gateway.REQUEST_TYPE(c.RequestType),
		AutoCapture:   c.AutoCapture,
		Status:        gateway.STATUS(c.Status),
		Metadata:      metadata,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}, nil
}

func (r *Repository) AddCardRequest(ctx context.Context, request gateway.CardRequest) (gateway.CardRequest, error) {
	db := r.DB(ctx)

	id, err := r.uidGenerator.NextID()
	if err != nil {
		return gateway.CardRequest{}, wrapDBError(err, "card_request")
	}
	request.ID = id

	var c cardRequestModel
	err = c.model(request)
	if err != nil {
		return gateway.CardRequest{}, wrapDBError(err, "card_request")
	}

	now := time.Now().Unix()
	c.CreatedAt = now
	c.UpdatedAt = now

	err = db.Create(&c).Error
	if err != nil {
		return gateway.CardRequest{}, wrapDBError(err, "card_request")
	}

	request.CreatedAt = c.CreatedAt
	request.UpdatedAt = c.UpdatedAt

	return request, nil
}

func (r *Repository) UpdateCardRequest(ctx context.Context, request gateway.CardRequest) (gateway.CardRequest, error) {
	db := r.DB(ctx)

	var c cardRequestModel
	err := c.model(request)
	if err != nil {
		return gateway.CardRequest{}, wrapDBError(err, "card_request")
	}

	c.UpdatedAt = time.Now().Unix()

	err = db.Model(&c).Updates(&c).Error
	if err != nil {
		return gateway.CardRequest{}, wrapDBError(err, "card_request")
	}

	request.UpdatedAt = c.UpdatedAt

	return request, nil
}
