package result

import (
	"errors"
	"github.com/shrinex/shrinex-core-backend/errx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
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
		body := &Result{Code: 200, Message: "OK", Data: resp}
		httpx.WriteJson(w, http.StatusOK, body)
		return
	}

	logx.WithContext(r.Context()).Errorf("[ERROR][API]: %+v ", err)

	// try to find ErrorEnvelope
	var e *errx.ErrorEnvelope
	if errors.As(err, &e) {
		httpx.WriteJson(w,
			http.StatusOK,
			&Result{
				Code:    e.Code,
				Message: e.Message,
			})
		return
	}

	httpx.WriteJson(w,
		http.StatusOK,
		&Result{
			Code:    errx.Unavailable.Code,
			Message: errx.Unavailable.Message,
		})
}
