package errx

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorEnvelope is a type that encapsulates Code and Message
type ErrorEnvelope struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	_                error = &ErrorEnvelope{}
	OK                     = &ErrorEnvelope{int(codes.OK), "OK"}
	Internal               = &ErrorEnvelope{int(codes.Internal), "Service unavailable, please try again later"}
	Unauthenticated        = &ErrorEnvelope{int(codes.Unauthenticated), "Unauthenticated"}
	PermissionDenied       = &ErrorEnvelope{int(codes.PermissionDenied), "Permission denied"}
)

// New creates a new instance of ErrorEnvelope
func New(code int, message string) error {
	return &ErrorEnvelope{
		Code:    code,
		Message: message,
	}
}

func (e *ErrorEnvelope) Error() string {
	return e.Message
}

// GRPCStatus convert ErrorEnvelope to gRPC Status
func (e *ErrorEnvelope) GRPCStatus() *status.Status {
	return status.New(codes.Code(e.Code), e.Message)
}
