package topup

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/common/uid"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/port/api"
)

type Service interface {
	AddTopUp(context.Context, uint64, api.AddTopUpRequest) (domain.TopUp, error)
}

type Repository interface {
	AddTopUp(context.Context, domain.TopUp) (domain.TopUp, error)

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

func (s *service) AddTopUp(ctx context.Context, userID uint64, req api.AddTopUpRequest) (domain.TopUp, error) {
	t, err := s.repo.TxnExec(ctx, func(ctxWithTxn context.Context) (interface{}, error) {
		topUpID, err := s.uidGenerator.NextID()
		if err != nil {
			return domain.TopUp{}, err
		}

		topUp := domain.TopUp{
			ID:              topUpID,
			CustomerID:      userID,
			PaymentMethodID: req.PaymentMethodID,
			Amount:          req.Amount,
			Currency:        req.Currency,
			Status:          domain.CREATED,
		}

		topUp, err = s.repo.AddTopUp(ctx, topUp)
		if err != nil {
			return domain.TopUp{}, err
		}

		return topUp, nil
	})

	if err != nil {
		return domain.TopUp{}, nil
	}

	return t.(domain.TopUp), nil
}
