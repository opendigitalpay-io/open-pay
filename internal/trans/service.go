package trans

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/common/uid"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type Service interface {
	AddTransfer(context.Context, domain.Order) (Transfer, error)
	UpdateTransfer(context.Context, Transfer) (Transfer, error)
}

type Repository interface {
	AddTransfer(context.Context, Transfer) (Transfer, error)
	UpdateTransfer(context.Context, Transfer) (Transfer, error)
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

func (s *service) AddTransfer(ctx context.Context, order domain.Order) (Transfer, error) {
	transferID, err := s.uidGenerator.NextID()
	if err != nil {
		return Transfer{}, err
	}
	transfer := Transfer{
		ID: transferID,
		OrderID: order.ID,
		CustomerID: order.CustomerID,
		SourceID: order.CustomerID, // FIXME: this should be the token from pay request
		DestinationID: order.MerchantID,
		Type: ORDER,
		Amount: order.Amount,
		Currency: order.Currency,
		Status: tcc.CREATED,
	}

	return s.repo.AddTransfer(ctx, transfer)
}

func (s *service) UpdateTransfer(ctx context.Context, transfer Transfer) (Transfer, error) {
	return s.repo.UpdateTransfer(ctx, transfer)
}
