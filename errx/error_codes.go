package errx

const (
	DataAccess     = 1000
	DataAccessDesc = "数据库访问异常"

	Validation     = 2000
	ValidationDesc = "参数校验失败"

	Regular     = 3000
	RegularDesc = "一般性异常"
)

func InvisibleCodes() []int {
	return []int{DataAccess}
}

func Visible(code int) bool {
	for _, c := range InvisibleCodes() {
		if c == code {
			return false
		}
	}
	return true
}
