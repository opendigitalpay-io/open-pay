package tcc

type Observer interface {
	OnTrySuccessCallback()
	OnTryFailCallback()
	OnCommitSuccessCallback()
	OnCommitFailCallback()
	OnCancelSuccessCallback()
	OnCancelFailCallback()
}
