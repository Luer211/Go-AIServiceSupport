package e

// 定义错误码
const (
	CodeSuccess = 0

	CodeInvalidParams = 40000
	CodeUnauthorized  = 40100
	CodeForbidden     = 40300
	CodeNotFound      = 40400
	CodeTooManyReq    = 42900

	CodeUserExists       = 40001
	CodeInvalidLogin     = 40002
	CodeTaskNotFound     = 40003
	CodeTaskSubmitFailed = 50001

	CodeInternalError = 50000
)

// 错误码对应的信息说明
var messages = map[int]string{
	CodeSuccess:          "success",
	CodeInvalidParams:    "invalid params",
	CodeUnauthorized:     "unauthorized",
	CodeForbidden:        "forbidden",
	CodeNotFound:         "not found",
	CodeTooManyReq:       "too many requests",
	CodeUserExists:       "user already exists",
	CodeInvalidLogin:     "invalid username or password",
	CodeTaskNotFound:     "task not found",
	CodeTaskSubmitFailed: "task submit failed",
	CodeInternalError:    "internal error",
}

// 错误码转 HTTP 状态码
func HTTPStatus(code int) int {
	switch code {
	case CodeInvalidParams, CodeUserExists, CodeInvalidLogin:
		return 400
	case CodeUnauthorized:
		return 401
	case CodeForbidden:
		return 403
	case CodeNotFound, CodeTaskNotFound:
		return 404
	case CodeTooManyReq:
		return 429
	case CodeInternalError, CodeTaskSubmitFailed:
		return 500
	default:
		return 500
	}
}

// 根据错误码返回错误信息
func Message(code int) string {
	if msg, ok := messages[code]; ok {
		return msg
	}
	return messages[CodeInternalError]
}
