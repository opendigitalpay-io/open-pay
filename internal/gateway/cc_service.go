package gateway

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/external/stripe"
	"github.com/opendigitalpay-io/open-pay/internal/common/uid"
	//"github.com/opendigitalpay-io/open-pay/internal/transtxn"
	"strconv"
)

type Service interface {
	CreateCardRequest(context.Context, CardRequestDTO) (CardRequest, error)
	Authorize(context.Context, CardRequest) (CardRequest, error)
	Capture(context.Context, CardRequest) (CardRequest, error)
}

type Repository interface {
	AddCardRequest(context.Context, CardRequest) (CardRequest, error)
	UpdateCardRequest(context.Context, CardRequest) (CardRequest, error)
}

type service struct {
	repo         Repository
	uidGenerator uid.Generator
	adapter      stripe.Adapter
}

func NewService(repo Repository, uidGenerator uid.Generator, adapter stripe.Adapter) Service {
	return &service{
		repo:         repo,
		uidGenerator: uidGenerator,
		adapter:      adapter,
	}
}

func (s *service) CreateCardRequest(ctx context.Context, cardRequestDTO CardRequestDTO) (CardRequest, error) {
	id, err := s.uidGenerator.NextID()
	if err != nil {
		return CardRequest{}, nil
	}

	cardRequest := CardRequest{
		ID:            id,
		TransferTxnID: cardRequestDTO.ID,
		GatewayToken:  cardRequestDTO.GatewayToken,
		// FIXME: better design for supporting multi-gateway
		Gateway:     STRIPE,
		Amount:      cardRequestDTO.Amount,
		Currency:    cardRequestDTO.Currency,
		RequestType: cardRequestDTO.RequestType,
		AutoCapture: cardRequestDTO.AutoCapture,
		Status:      CREATED,
		Metadata:    cardRequestDTO.Metadata,
	}

	cardRequest, err = s.repo.AddCardRequest(ctx, cardRequest)
	if err != nil {
		return CardRequest{}, err
	}

	return cardRequest, nil
}

func (s *service) Authorize(ctx context.Context, request CardRequest) (CardRequest, error) {
	requestID, err := s.uidGenerator.NextID()
	if err != nil {
		return CardRequest{}, err
	}

	request.RequestID = requestID
	stripeChargeRequest := stripe.StripeChargeRequest{
		Amount:    request.Amount,
		Currency:  request.Currency,
		Source:    request.GatewayToken,
		RequestID: strconv.FormatUint(request.RequestID, 10),
	}
	gatewayTxnID, err := s.adapter.Authorize(stripeChargeRequest)
	if err != nil {
		request.Status = FAILED
		_, dbErr := s.repo.UpdateCardRequest(ctx, request)
		// TODO: how do we solve this nested error?
		if dbErr != nil {
			return CardRequest{}, dbErr
		}

		return CardRequest{}, err
	}

	request.GatewayTxnID = gatewayTxnID
	request.Status = COMPLETED
	request, err = s.repo.UpdateCardRequest(ctx, request)
	if err != nil {
		return CardRequest{}, err
	}

	return request, nil
}

func (s *service) Capture(ctx context.Context, request CardRequest) (CardRequest, error) {
	requestID, err := s.uidGenerator.NextID()
	if err != nil {
		return CardRequest{}, err
	}
	request.RequestID = requestID
	stripeCaptureRequest := stripe.StripeCaptureRequest{
		ChargeID:  request.GatewayTxnID,
		RequestID: strconv.FormatUint(request.RequestID, 10),
	}
	err = s.adapter.Capture(stripeCaptureRequest)
	if err != nil {
		request.Status = FAILED
		_, dbErr := s.repo.UpdateCardRequest(ctx, request)
		// TODO: how do we solve this nested error?
		if dbErr != nil {
			return CardRequest{}, dbErr
		}
		return CardRequest{}, err
	}
	request.Status = COMPLETED
	request, err = s.repo.UpdateCardRequest(ctx, request)
	if err != nil {
		return CardRequest{}, err
	}

	return request, nil
}
