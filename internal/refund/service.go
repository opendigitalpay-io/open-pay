package refund

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/common/uid"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
)

type Service interface {
	AddRefund(context.Context, uint64) (domain.Refund, error)
}

type Repository interface {
	AddRefund(context.Context, domain.Refund) (domain.Refund, error)

	GetOrder(context.Context, uint64) (domain.Order, error)

	TxnExec(context.Context, func(context.Context) (interface{}, error)) (interface{}, error)
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

func (s *service) AddRefund(ctx context.Context, orderID uint64) (domain.Refund, error) {
	t, err := s.repo.TxnExec(ctx, func(ctx context.Context) (interface{}, error) {
		refundID, err := s.uidGenerator.NextID()
		if err != nil {
			return domain.Refund{}, err
		}

		order, err := s.repo.GetOrder(ctx, orderID)

		if err != nil {
			return domain.Refund{}, err
		}

		refund := domain.Refund{
			ID:          refundID,
			OrderID:     order.ID,
			Amount:      order.Amount,
			Status:      domain.CREATED,
			RefundCount: 0,
		}

		refund, err = s.repo.AddRefund(ctx, refund)
		if err != nil {
			return domain.Refund{}, err
		}

		return refund, nil
	})

	if err != nil {
		return domain.Refund{}, nil
	}

	return t.(domain.Refund), nil
}
