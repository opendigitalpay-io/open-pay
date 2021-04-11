package tcc

import (
	"context"
)

type Strategy interface {
	Try(context.Context) error
	Commit(context.Context) error
	Cancel(context.Context) error
	GetStatus() STATUS
	AddObserver(Observer)
}
