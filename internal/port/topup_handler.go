package port

import (
	"github.com/gin-gonic/gin"
	"github.com/opendigitalpay-io/open-pay/internal/port/api"
)

func (h *HTTPServer) AddTopUp() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var topUpUriParam api.AddTopUpUriParameter
		var topUpRequest api.AddTopUpRequest

		if err := ctx.ShouldBindUri(&topUpUriParam); err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		if err := ctx.ShouldBindJSON(&topUpRequest); err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		topUp, err := h.topUpService.AddTopUp(ctx, topUpUriParam.ID, topUpRequest)

		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		addTopupResp := api.AddTopUpResponse{
			ID: topUp.ID,
		}

		h.RespondWithOK(ctx, addTopupResp)
	}
}
