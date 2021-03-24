package port

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/opendigitalpay-io/open-pay/internal/common/errorz"
	"github.com/opendigitalpay-io/open-pay/internal/common/server"
	"github.com/opendigitalpay-io/open-pay/internal/common/uid"
	"github.com/opendigitalpay-io/open-pay/internal/refund"
	"github.com/opendigitalpay-io/open-pay/internal/storage"
	"github.com/opendigitalpay-io/open-pay/internal/topup"
)

type HTTPServer struct {
	topupService  topup.Service
	refundService refund.Service
	repo          *storage.Repository // FIXME: this should be removed
	uidGenerator  uid.Generator       // FIXME: this should be removed
}

func NewHTTPServer(topupService topup.Service, refundService refund.Service, repo *storage.Repository, uidGenerator uid.Generator) *HTTPServer {
	return &HTTPServer{
		topupService:  topupService,
		refundService: refundService,
		repo:          repo,
		uidGenerator:  uidGenerator,
	}
}

func (*HTTPServer) RespondWithOK(ctx *gin.Context, resp interface{}) {
	server.OK(ctx, resp)
}

func (*HTTPServer) RespondWithError(ctx *gin.Context, err error) {
	var ves validator.ValidationErrors
	if errors.As(err, &ves) {
		server.BadRequest(ctx, errorz.NewValidationError(ves), err)
		return
	}

	var syne *json.SyntaxError
	if errors.As(err, &syne) {
		server.BadRequest(ctx, errorz.NewInvalidJSONError(syne), err)
		return
	}
	//
	//var nfe storage.NotFoundError
	//if errors.As(err, &nfe) {
	//	server.NotFound(ctx, errorz.NewNotFoundError(nfe), err)
	//	return
	//}
	//
	//var dee storage.DuplicatedEntryError
	//if errors.As(err, &dee) {
	//	server.BadRequest(ctx, errorz.NewInvalidValueError(dee), err)
	//	return
	//}
	//
	//var ideme domain.IdemError
	//if errors.As(err, &ideme) {
	//	server.BadRequest(ctx, errorz.NewIdemKeyError(err), err)
	//	return
	//}

	// fallback error handling
	server.InternalError(ctx, errorz.NewInternalError(err), err)
}
