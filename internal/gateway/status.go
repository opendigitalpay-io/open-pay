package gateway

type STATUS string

const (
	CREATED   STATUS = "CREATED"
	COMPLETED STATUS = "COMPLETED"
	FAILED    STATUS = "FAILED"
)

var status = [...]string{
	"CREATED",
	"COMPLETED",
	"FAILED",
}

func (s *STATUS) String() string {
	x := string(*s)
	for _, v := range status {
		if v == x {
			return x
		}
	}

	return ""
}
