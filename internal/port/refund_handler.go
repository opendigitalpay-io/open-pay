package port

import (
	"github.com/gin-gonic/gin"
	"github.com/opendigitalpay-io/open-pay/internal/port/api"
)

func (h *HTTPServer) AddRefund() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var refundUriParam api.AddRefundUriParameter

		if err := ctx.ShouldBindUri(&refundUriParam); err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		refund, err := h.refundService.AddRefund(ctx, refundUriParam.OrderID)

		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		addRefundResp := api.AddRefundResponse{
			ID: refund.ID,
		}

		h.RespondWithOK(ctx, addRefundResp)
	}

}
