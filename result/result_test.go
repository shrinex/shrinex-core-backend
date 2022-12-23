package result

import (
	"errors"
	"fmt"
	"github.com/shrinex/shrinex-core-backend/errx"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/jsonx"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompute_WhenNoError_Render200OK(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()
	resp := map[string]any{"key": "value"}
	Compute(r, w, resp, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var ret Result
	buf, err := ioutil.ReadAll(w.Body)
	if assert.NoError(t, err) {
		_ = jsonx.Unmarshal(buf, &ret)
		assert.Equal(t, 200, ret.Code)
		assert.Equal(t, "OK", ret.Message)
		assert.NotNil(t, ret.Data)
		assert.IsType(t, resp, ret.Data)
		assert.Equal(t, "value", (ret.Data.(map[string]any))["key"])
	}
}

func TestCompute_WhenErrorIsNotErrorEnvelope_ReturnsUnavailable(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()
	Compute(r, w, nil, errors.New("A fresh error"))

	var ret Result
	buf, err := ioutil.ReadAll(w.Body)
	if assert.NoError(t, err) {
		_ = jsonx.Unmarshal(buf, &ret)
		assert.Equal(t, errx.Unavailable.Code, ret.Code)
		assert.Equal(t, errx.Unavailable.Message, ret.Message)
		assert.Nil(t, ret.Data)
	}
}

func TestCompute_WhenErrorIsErrorEnvelope_ReturnsItsCodeAndMessage(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()
	Compute(r, w, nil, fmt.Errorf("%w: timeout occurs",
		errx.New(1024, "Timeout")))

	var ret Result
	buf, err := ioutil.ReadAll(w.Body)
	if assert.NoError(t, err) {
		_ = jsonx.Unmarshal(buf, &ret)
		assert.Equal(t, 1024, ret.Code)
		assert.Equal(t, "Timeout", ret.Message)
		assert.Nil(t, ret.Data)
	}
}
