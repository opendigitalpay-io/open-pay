package tcc

import "context"

type Observer interface {
	OnTrySuccessCallback(context.Context)
	OnTryFailCallback(context.Context)
	OnCommitSuccessCallback(context.Context)
	OnCommitFailCallback(context.Context)
	OnCancelSuccessCallback(context.Context)
	OnCancelFailCallback(context.Context)
}
