package errx

import (
	"strconv"
	"sync"
)

const (
	// Regular 是一般性异常，客
	// 户端不关心错误码时使用
	Regular = 10000

	// DataAccess 是数据访问异常，
	// 这种类型的异常一般不会透给客户端
	DataAccess = 10001

	// Validation 是参数校验失败
	// 产生的异常
	Validation = 10002
)

func init() {
	Register(DataAccess)
}

var (
	invisibleCodesMu sync.RWMutex
	invisibleCodes   = make([]int, 0, 8)
)

// Register registers specified invisible code
func Register(code int) {
	invisibleCodesMu.Lock()
	defer invisibleCodesMu.Unlock()
	if !visibleLocked(code) {
		panic("errx: Register called twice for code " + strconv.Itoa(code))
	}
	invisibleCodes = append(invisibleCodes, code)
}

// InvisibleCodes returns an immutable copy of all invisible codes
func InvisibleCodes() []int {
	invisibleCodesMu.RLock()
	defer invisibleCodesMu.RUnlock()
	var result = make([]int, 0, len(invisibleCodes))
	for _, code := range invisibleCodes {
		result = append(result, code)
	}
	return result
}

// Visible tests whether a specified code is visible
func Visible(code int) bool {
	invisibleCodesMu.RLock()
	defer invisibleCodesMu.RUnlock()
	return visibleLocked(code)
}

func visibleLocked(code int) bool {
	for _, c := range invisibleCodes {
		if c == code {
			return false
		}
	}
	return true
}
