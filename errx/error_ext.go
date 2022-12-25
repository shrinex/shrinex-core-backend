package errx

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
)

// Is test if err contains specified error code
func Is(err error, code int, codes ...int) bool {
	codes = append(codes, code)

	if err == nil {
		return intSlice(codes).Contains(OK.Code)
	}

	// try ErrorEnvelope first
	var e *ErrorEnvelope
	if errors.As(err, &e) {
		return intSlice(codes).Contains(e.Code)
	} else {
		// or maybe its gRPC error
		for err != nil {
			if s, ok := status.FromError(err); ok {
				return intSlice(codes).Contains(int(s.Code()))
			}
			err = errors.Unwrap(err)
		}
	}

	// not code error
	return false
}

type intSlice []int

func (is intSlice) Contains(e int) bool {
	for i := 0; i < len(is); i++ {
		if e == is[i] {
			return true
		}
	}
	return false
}
