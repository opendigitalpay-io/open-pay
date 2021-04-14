package port

import (
	"github.com/gin-gonic/gin"
	"github.com/opendigitalpay-io/open-pay/internal/port/api"
)

func (h *HTTPServer) AddOrder() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req api.AddOrderRequest

		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		order, err := h.orderService.AddOrder(ctx, req)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		resp := api.AddOrderResponse{
			OrderID: order.ID,
		}

		h.RespondWithOK(ctx, resp)
	}
}

func (h *HTTPServer) GetOrder() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var uriParam api.GetOrderURIParameter
		err := ctx.ShouldBindUri(&uriParam)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		order, err := h.orderService.GetOrder(ctx, uriParam.ID)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		lineItems, ok := order.Metadata["lineItems"].([]api.LineItem)
		if ok {
			delete(order.Metadata, "lineItems")
		}

		resp := api.GerOrderResponse{
			ID:            order.ID,
			CustomerID:    order.CustomerID,
			BusinessID:    order.MerchantID,
			Amount:        order.Amount,
			Currency:      order.Currency,
			ReferenceID:   order.ReferenceID,
			CustomerEmail: order.CustomerEmail,
			LineItems:     lineItems,
			Metadata:      order.Metadata,
			CreatedAt:     order.CreatedAt,
		}

		h.RespondWithOK(ctx, resp)
	}
}

func (h *HTTPServer) PayOrder() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var uriParam api.PayOrderURIParameter
		err := ctx.ShouldBindUri(&uriParam)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		var req api.PayOrderRequest

		err = ctx.ShouldBindJSON(&req)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		err = h.orderPayService.PayOrder(ctx, uriParam.ID, req)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		h.RespondWithOK(ctx, nil)
	}
}