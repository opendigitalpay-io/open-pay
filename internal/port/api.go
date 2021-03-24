package port

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ServerInterface interface {
	// GET /health
	GetHealthStatus() func(*gin.Context)

	// POST /v1/user/{id}/topup
	AddTopup() func(*gin.Context)

	// POST /v1/order/:id/refund
	AddRefund() func(ctx *gin.Context)
}

func HandlerFromMux(si ServerInterface, e *gin.Engine) http.Handler {
	// Healthz
	e.GET("/health", si.GetHealthStatus())

	v1 := e.Group("/v1")
	{
		// topup
		v1.POST("/user/:id/topup", si.AddTopup())

		// refund
		v1.POST("/order/:id/refund", si.AddRefund())
	}

	return e
}
