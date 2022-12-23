package errx

// ErrorEnvelope is a type that encapsulates Code and Message
type ErrorEnvelope struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	_           error = &ErrorEnvelope{}
	Unavailable       = &ErrorEnvelope{-1, "Service unavailable, please try again later"}
)

func (e *ErrorEnvelope) Error() string {
	return e.Message
}

// New creates a new instance of ErrorEnvelope
func New(code int, message string) error {
	return &ErrorEnvelope{
		Code:    code,
		Message: message,
	}
}
