package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/opendigitalpay-io/open-pay/external/balance"
	"github.com/opendigitalpay-io/open-pay/external/stripe"
	"github.com/opendigitalpay-io/open-pay/internal/common/server"
	"github.com/opendigitalpay-io/open-pay/internal/common/uid"
	"github.com/opendigitalpay-io/open-pay/internal/factory"
	"github.com/opendigitalpay-io/open-pay/internal/gateway"
	"github.com/opendigitalpay-io/open-pay/internal/order"
	"github.com/opendigitalpay-io/open-pay/internal/pay"
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
	stripeAdapter := stripe.NewAdapter()

	gatewayService := gateway.NewService(repository, uidGenerator, stripeAdapter)
	transferTxnService := transtxn.NewService(repository, balanceAdapter, gatewayService, uidGenerator)
	transferService := trans.NewService(repository, uidGenerator)
	orderService := order.NewService(repository, uidGenerator)

	transferTxnStrategyFactory := factory.NewTransferTxnStrategyFactory(transferTxnService)
	transferStrategyFactory := factory.NewTransferStrategyFactory(transferService, transferTxnStrategyFactory)
	orderStrategyFactory := factory.NewOrderStrategyFactory(orderService, transferStrategyFactory)

	payService := pay.NewService(orderStrategyFactory)
	topUpService := topup.NewService(repository, uidGenerator)
	refundService := refund.NewService(repository, uidGenerator)

	server.RunHTTPServer(func(engine *gin.Engine) http.Handler {
		return port.HandlerFromMux(
			port.NewHTTPServer(orderService, payService, topUpService, refundService),
			engine,
		)
	})
}
