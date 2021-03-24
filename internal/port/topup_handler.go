package port

import (
	"github.com/gin-gonic/gin"
	"github.com/opendigitalpay-io/open-pay/internal/port/api"
)

func (h *HTTPServer) AddTopup() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var topupUriParam api.AddTopupUriParameter
		var topupRequest api.AddTopupRequest

		if err := ctx.ShouldBindUri(&topupUriParam); err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		if err := ctx.ShouldBindJSON(&topupRequest); err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		topup, err := h.topupService.AddTopup(ctx, topupUriParam.ID, topupRequest)

		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		addTopupResp := api.AddTopupResponse{
			ID: topup.ID,
		}

		h.RespondWithOK(ctx, addTopupResp)
	}
}
