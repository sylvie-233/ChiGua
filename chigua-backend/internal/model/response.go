package model

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// ResponseCode 响应码枚举
type ResponseCode int

const (
	// 成功
	Success ResponseCode = 200

	// 客户端错误
	BadRequest          ResponseCode = 400
	Unauthorized        ResponseCode = 401
	Forbidden           ResponseCode = 403
	NotFound            ResponseCode = 404
	MethodNotAllowed    ResponseCode = 405
	RequestTimeout      ResponseCode = 408
	Conflict            ResponseCode = 409
	UnprocessableEntity ResponseCode = 422

	// 服务器错误
	InternalServerError ResponseCode = 500
	ServiceUnavailable  ResponseCode = 503
	GatewayTimeout      ResponseCode = 504

	// 业务错误
	InvalidParams    ResponseCode = 1000
	UserExists       ResponseCode = 1001
	UserNotExists    ResponseCode = 1002
	PasswordError    ResponseCode = 1003
	TokenExpired     ResponseCode = 1004
	TokenInvalid     ResponseCode = 1005
	ArticleNotFound  ResponseCode = 2000
	CategoryNotFound ResponseCode = 2001
	TagNotFound      ResponseCode = 2002
)

// GetMsg 获取响应码对应的消息
func (c ResponseCode) GetMsg() string {
	switch c {
	case Success:
		return "成功"
	case BadRequest:
		return "请求参数错误"
	case Unauthorized:
		return "未授权"
	case Forbidden:
		return "禁止访问"
	case NotFound:
		return "资源不存在"
	case MethodNotAllowed:
		return "方法不允许"
	case RequestTimeout:
		return "请求超时"
	case Conflict:
		return "资源冲突"
	case UnprocessableEntity:
		return "无法处理的实体"
	case InternalServerError:
		return "服务器内部错误"
	case ServiceUnavailable:
		return "服务不可用"
	case GatewayTimeout:
		return "网关超时"
	case InvalidParams:
		return "参数无效"
	case UserExists:
		return "用户已存在"
	case UserNotExists:
		return "用户不存在"
	case PasswordError:
		return "密码错误"
	case TokenExpired:
		return "令牌已过期"
	case TokenInvalid:
		return "令牌无效"
	case ArticleNotFound:
		return "文章不存在"
	case CategoryNotFound:
		return "分类不存在"
	case TagNotFound:
		return "标签不存在"
	default:
		return "未知错误"
	}
}

// NewResponse 创建新的响应
func NewResponse(code ResponseCode, data interface{}) Response {
	return Response{
		Code: int(code),
		Msg:  code.GetMsg(),
		Data: data,
	}
}

// SuccessResponse 创建成功响应
func SuccessResponse(data interface{}) Response {
	return NewResponse(Success, data)
}

// ErrorResponse 创建错误响应
func ErrorResponse(code ResponseCode) Response {
	return NewResponse(code, nil)
}
