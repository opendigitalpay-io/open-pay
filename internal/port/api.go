package port

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ServerInterface interface {
	// GET /health
	GetHealthStatus() func(*gin.Context)

	// GET /v1/order/:id
	GetOrder() func(*gin.Context)
	// POST /v1/order/
	AddOrder() func(*gin.Context)
	// POST /v1/order/:id/refund
	AddRefund() func(*gin.Context)

	// POST /v1/user/{id}/topup
	AddTopUp() func(*gin.Context)
}

func HandlerFromMux(si ServerInterface, e *gin.Engine) http.Handler {
	// Healthz
	e.GET("/health", si.GetHealthStatus())

	v1 := e.Group("/v1")
	{
		// order
		o := v1.Group("/order")
		{
			o.GET("/:id", si.AddOrder())
			o.POST("", si.AddOrder())
			o.POST("/:id/refund", si.AddRefund())
		}

		// topup
		v1.POST("/user/:id/topup", si.AddTopUp())

	}

	return e
}
