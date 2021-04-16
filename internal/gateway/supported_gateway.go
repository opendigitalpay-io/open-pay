package gateway

type SUPPORTED_GATEWAY string

const (
	STRIPE SUPPORTED_GATEWAY = "STRIPE"
)

var supported_gateway = [...]string{
	"STRIPE",
}

func (s *SUPPORTED_GATEWAY) String() string {
	x := string(*s)
	for _, v := range supported_gateway {
		if v == x {
			return x
		}
	}

	return ""
}
