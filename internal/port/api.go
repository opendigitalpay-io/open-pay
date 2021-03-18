package port

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ServerInterface interface {
	// GET /health
	GetHealthStatus() func(*gin.Context)

}

func HandlerFromMux(si ServerInterface, e *gin.Engine) http.Handler {
	// Healthz
	e.GET("/health", si.GetHealthStatus())

	return e
}
