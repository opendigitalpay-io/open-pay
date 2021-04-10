package tcc

import (
	"context"
	"github.com/opendigitalpay-io/open-pay/internal/domain"
)

type Strategy interface {
	Try(context.Context) error
	Commit(context.Context) error
	Cancel(context.Context) error
	GetStatus() domain.STATUS
}
