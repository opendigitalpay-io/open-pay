package domain

type IdemError struct {
	What string
}

func (e IdemError) Error() string {
	return e.What
}
