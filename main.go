package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/opendigitalpay-io/open-pay/external/balance"
	"github.com/opendigitalpay-io/open-pay/internal/common/server"
	"github.com/opendigitalpay-io/open-pay/internal/common/uid"
	"github.com/opendigitalpay-io/open-pay/internal/order"
	"github.com/opendigitalpay-io/open-pay/internal/port"
	"github.com/opendigitalpay-io/open-pay/internal/refund"
	"github.com/opendigitalpay-io/open-pay/internal/storage"
	"github.com/opendigitalpay-io/open-pay/internal/topup"
	"github.com/opendigitalpay-io/open-pay/internal/trans"
	"github.com/opendigitalpay-io/open-pay/internal/transtxn"
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

	balanceAdapter := balance.NewAdapter()

	transferTxnService := transtxn.NewService(repository, balanceAdapter, uidGenerator)
	transferService := trans.NewService(repository, uidGenerator)
	orderService := order.NewService(repository, uidGenerator)

	transferTxnStrategyFactory := transtxn.NewStrategyFactory(transferTxnService)
	transferStrategyFactory := trans.NewStrategyFactory(transferService, transferTxnStrategyFactory)
	orderStrategyFactory := order.NewStrategyFactory(orderService, transferStrategyFactory)

	orderPayService := order.NewPayService(orderStrategyFactory)
	topUpService := topup.NewService(repository, uidGenerator)
	refundService := refund.NewService(repository, uidGenerator)

	server.RunHTTPServer(func(engine *gin.Engine) http.Handler {
		return port.HandlerFromMux(
			port.NewHTTPServer(orderService, orderPayService, topUpService, refundService),
			engine,
		)
	})
}
