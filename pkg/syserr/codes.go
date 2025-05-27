package syserr

type Code string

// System error codes.
const (
	InternalCode        Code = "internal"
	InvalidArgumentCode Code = "invalid_argument"
)
