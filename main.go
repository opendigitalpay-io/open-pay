package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/opendigitalpay-io/open-pay/internal/common/server"
	"github.com/opendigitalpay-io/open-pay/internal/common/uid"
	"github.com/opendigitalpay-io/open-pay/internal/port"
	"github.com/opendigitalpay-io/open-pay/internal/storage"
	"net/http"
)

func main() {
	ctx := context.Background()
	repository, err := storage.NewRepository(ctx, &storage.Config{})
	if err != nil {
		panic(err)
	}

	uidGenerator, err := uid.NewGenerator(ctx)

	if err != nil {
		panic(err)
	}

	server.RunHTTPServer(func(engine *gin.Engine) http.Handler {
		return port.HandlerFromMux(
			port.NewHTTPServer(repository, uidGenerator),
			engine,
		)
	})
}
