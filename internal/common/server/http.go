package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/opendigitalpay-io/open-pay/internal/common/validator"
	"net/http"
)

func RunHTTPServer(createHandler func(engine *gin.Engine) http.Handler) {
	router := gin.Default()

	binding.Validator = validator.NewDefaultValidator()

	createHandler(router)

	http.ListenAndServe(":8182", router)
}
