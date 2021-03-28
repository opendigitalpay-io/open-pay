package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/opendigitalpay-io/open-pay/internal/common/server"
	"github.com/opendigitalpay-io/open-pay/internal/common/uid"
	"github.com/opendigitalpay-io/open-pay/internal/order"
	"github.com/opendigitalpay-io/open-pay/internal/port"
	"github.com/opendigitalpay-io/open-pay/internal/refund"
	"github.com/opendigitalpay-io/open-pay/internal/storage"
	"github.com/opendigitalpay-io/open-pay/internal/topup"
	"net/http"
)

func main() {
	ctx := context.Background()
	uidGenerator, err := uid.NewGenerator(ctx)

	if err != nil {
		panic(err)
	}

	repository, err := storage.NewRepository(ctx, &storage.Config{}, uidGenerator)
	if err != nil {
		panic(err)
	}

	orserService := order.NewService(repository, uidGenerator)
	topupService := topup.NewService(repository, uidGenerator)
	refundService := refund.NewService(repository, uidGenerator)

	server.RunHTTPServer(func(engine *gin.Engine) http.Handler {
		return port.HandlerFromMux(
			port.NewHTTPServer(orserService, topupService, refundService),
			engine,
		)
	})
}
