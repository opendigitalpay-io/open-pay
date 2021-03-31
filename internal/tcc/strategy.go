package tcc

import "github.com/opendigitalpay-io/open-pay/internal/domain"

type Strategy interface {
	Try() error
	Commit() error
	Cancel() error
	GetStatus() domain.STATUS


}
