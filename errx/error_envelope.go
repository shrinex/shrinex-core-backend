package errx

import (
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

// ErrorEnvelope is a type that encapsulates Code and Message
type ErrorEnvelope struct {
	Code    int
	Message string
	cause   error
}

type Options func(*ErrorEnvelope)

var (
	_  error = &ErrorEnvelope{}
	OK       = &ErrorEnvelope{
		Code:    int(codes.OK),
		Message: "OK",
	}
	Internal = &ErrorEnvelope{
		Code:    int(codes.Internal),
		Message: "Service unavailable, please try again later",
	}
	Unauthenticated = &ErrorEnvelope{
		Code:    int(codes.Unauthenticated),
		Message: "Unauthenticated",
	}
	PermissionDenied = &ErrorEnvelope{
		Code:    int(codes.PermissionDenied),
		Message: "Permission denied",
	}
	TokenExpired = &ErrorEnvelope{
		Code:    99999,
		Message: "Session expired, please re-login",
	}

	errStackTrace = errors.New("stack trace")
)

// New creates a new instance of ErrorEnvelope
func New(code int, message string, opts ...Options) error {
	result := &ErrorEnvelope{
		Code:    code,
		Message: message,
	}

	for _, opt := range opts {
		opt(result)
	}

	if result.cause == nil {
		result.cause = errors.WithStack(errStackTrace)
	}

	return result
}

// NewRegular creates a new ErrorEnvelope with Regular error code
func NewRegular(message string, opts ...Options) error {
	return New(Regular, message, opts...)
}

// NewValidation creates a new ErrorEnvelope with Validation error code
func NewValidation(message string, opts ...Options) error {
	return New(Validation, message, opts...)
}

// NewDataAccess creates a new ErrorEnvelope with DataAccess error code
func NewDataAccess(message string, opts ...Options) error {
	return New(DataAccess, message, opts...)
}

func (e *ErrorEnvelope) Error() string {
	if e.cause == nil || e.cause == errStackTrace {
		return e.Message
	}

	return e.Message + ": " + e.cause.Error()
}

func (e *ErrorEnvelope) Unwrap() error {
	return e.cause
}

func (e *ErrorEnvelope) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", e.cause)
			_, _ = io.WriteString(s, e.Message)
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, e.Error())
	}
}

// GRPCStatus convert ErrorEnvelope to gRPC Status
func (e *ErrorEnvelope) GRPCStatus() *status.Status {
	return status.New(codes.Code(e.Code), e.Message)
}

// WithCause customize ErrorEnvelope with a root cause
func WithCause(cause error) Options {
	return func(e *ErrorEnvelope) {
		e.cause = errors.WithStack(cause)
	}
}
