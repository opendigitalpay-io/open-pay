package domain

type STATUS string

const (
	CREATED         STATUS = "CREATED"
	PREPARE_STARTED STATUS = "PREPARE_STARTED"
	PREPARE_FAILED  STATUS = "PREPARE_FAILED"
	PREPARED        STATUS = "PREPARED"

	CONFIRM_STARTED STATUS = "CONFIRM_STARTED"
	CONFIRM_FAILED  STATUS = "CONFIRM_FAILED"
	CONFIRMED       STATUS = "CONFIRMED"

	CANCEL_STARTED STATUS = "CANCEL_STARTED"
	CANCEL_FAILED  STATUS = "CANCEL_FAILED"
	CANCELLED      STATUS = "CANCELLED"
)

var status = [...]string{
	"CREATED",
	"PREPARE_STARTED",
	"PREPARE_FAILED",
	"PREPARED",
	"CONFIRM_STARTED",
	"CONFIRM_FAILED",
	"CONFIRMED",
	"CANCEL_STARTED",
	"CANCEL_FAILED",
	"CANCELLED",
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
