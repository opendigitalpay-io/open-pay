package order

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/common/uid"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/port/api"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type Service interface {
	AddOrder(context.Context, api.AddOrderRequest) (domain.Order, error)
	GetOrder(context.Context, uint64) (domain.Order, error)
	UpdateOrder(context.Context, domain.Order) (domain.Order, error)
}

type Repository interface {
	AddOrder(context.Context, domain.Order) (domain.Order, error)
	GetOrder(context.Context, uint64) (domain.Order, error)
	UpdateOrder(context.Context, domain.Order) (domain.Order, error)
}

type service struct {
	repo         Repository
	uidGenerator uid.Generator
}

func NewService(repo Repository, uidGenerator uid.Generator) Service {
	return &service{
		repo:         repo,
		uidGenerator: uidGenerator,
	}
}

func (s *service) AddOrder(ctx context.Context, request api.AddOrderRequest) (domain.Order, error) {
	orderID, err := s.uidGenerator.NextID()
	if err != nil {
		return domain.Order{}, err
	}

	metadata := request.Metadata
	metadata["lineItems"] = request.LineItems

	order := domain.Order{
		ID:            orderID,
		CustomerID:    request.CustomerID,
		MerchantID:    request.BusinessID,
		Amount:        request.Amount,
		Currency:      request.Currency,
		ReferenceID:   request.ReferenceID,
		CustomerEmail: request.CustomerEmail,
		Status:        tcc.CREATED,
		Mode:          domain.INDIRECT,
		Metadata:      metadata,
	}

	order, err = s.repo.AddOrder(ctx, order)
	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (s *service) GetOrder(ctx context.Context, orderID uint64) (domain.Order, error) {
	return s.repo.GetOrder(ctx, orderID)
}

func (s *service) UpdateOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	return s.repo.UpdateOrder(ctx, order)
}
