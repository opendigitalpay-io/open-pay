package gateway

type REQUEST_TYPE string

const (
	AUTHORIZE REQUEST_TYPE = "AUTHORIZE"
	CAPTURE   REQUEST_TYPE = "CAPTURE"
	REFUND    REQUEST_TYPE = "REFUND"
)

var request_type = [...]string{
	"CREATED",
	"COMPLETED",
	"FAILED",
}

func (s *REQUEST_TYPE) String() string {
	x := string(*s)
	for _, v := range request_type {
		if v == x {
			return x
		}
	}

	return ""
}
