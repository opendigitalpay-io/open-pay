package topup

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/common/uid"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/port/api"
)

type Service interface {
	AddTopup(context.Context, uint64, api.AddTopupRequest) (domain.Topup, error)
}

type Repository interface {
	AddTopup(context.Context, domain.Topup) (domain.Topup, error)

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

func (s *service) AddTopup(ctx context.Context, userID uint64, req api.AddTopupRequest) (domain.Topup, error) {
	t, err := s.repo.TxnExec(ctx, func(ctxWithTxn context.Context) (interface{}, error) {
		topupID, err := s.uidGenerator.NextID()
		if err != nil {
			return domain.Topup{}, err
		}

		topup := domain.Topup{
			ID:              topupID,
			CustomerID:      userID,
			PaymentMethodID: req.PaymentMethodID,
			Amount:          req.Amount,
			Currency:        req.Currency,
			Status:          domain.CREATED,
		}

		topup, err = s.repo.AddTopup(ctx, topup)
		if err != nil {
			return domain.Topup{}, err
		}

		return topup, nil
	})

	if err != nil {
		return domain.Topup{}, nil
	}

	return t.(domain.Topup), nil
}
