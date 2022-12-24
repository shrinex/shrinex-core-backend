package result

import (
	"errors"
	"github.com/shrinex/shrinex-core-backend/errx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
	"net/http"
)

// Result is a type that represents common HTTP response
type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// Compute setup Result for writing HTTP response
func Compute(r *http.Request, w http.ResponseWriter, resp any, err error) {
	if err == nil {
		httpx.OkJsonCtx(r.Context(), w, &Result{
			Code:    errx.OK.Code,
			Message: errx.OK.Message,
			Data:    resp,
		})
		return
	}

	logx.WithContext(r.Context()).Errorf("[ERROR]: %+v ", err)

	// try ErrorEnvelope first
	var e *errx.ErrorEnvelope
	if errors.As(err, &e) {
		if errx.Visible(e.Code) {
			httpx.OkJsonCtx(r.Context(), w, &Result{
				Code:    e.Code,
				Message: e.Message,
			})
			return
		}
	} else {
		// or maybe its gRPC error
		for err != nil {
			if s, ok := status.FromError(err); ok {
				if errx.Visible(int(s.Code())) {
					httpx.OkJsonCtx(r.Context(), w, &Result{
						Code:    int(s.Code()),
						Message: s.Message(),
					})
					return
				} else {
					break
				}
			}
			err = errors.Unwrap(err)
		}
	}

	// ok, let's say it's an Internal error
	httpx.OkJsonCtx(r.Context(), w, &Result{
		Code:    errx.Internal.Code,
		Message: errx.Internal.Message,
	})
}
