package errors

// 业务错误码定义
// 格式：4xx = 客户端错误，5xx = 服务端错误

// 认证错误 401xx
var (
	ErrUnauthorized   = New(40101, "未登录")
	ErrTokenExpired   = New(40102, "Token 已过期")
	ErrTokenInvalid   = New(40103, "Token 无效")
	ErrCodeWrong      = New(40104, "验证码错误")
	ErrCodeExpired    = New(40105, "验证码已过期")
	ErrCodeFrequent   = New(40106, "验证码发送过于频繁，请60秒后重试")
)

// 用户错误 402xx
var (
	ErrPhoneExists    = New(40201, "手机号已注册")
	ErrNicknameExists = New(40202, "昵称已被使用")
	ErrUserNotFound   = New(40203, "用户不存在")
	ErrPasswordWrong  = New(40204, "密码错误")
	ErrUserDisabled   = New(40205, "账号已被禁用")
	ErrEmailExists    = New(40206, "邮箱已注册")
)

// 帖子错误 403xx
var (
	ErrPostNotFound   = New(40301, "帖子不存在")
	ErrPostForbidden  = New(40302, "无权操作此帖子")
	ErrTitleTooLong   = New(40303, "标题长度不能超过200字符")
	ErrContentEmpty   = New(40304, "内容不能为空")
	ErrTypeInvalid    = New(40305, "无效的帖子类型")
	ErrContentBlocked = New(40306, "内容包含违规词汇")
)

// 标签错误 404xx
var (
	ErrTagNotFound    = New(40401, "标签不存在")
)

// 通用错误 400xx
var (
	ErrParamInvalid   = New(40001, "参数错误")
	ErrParamMissing   = New(40002, "缺少必填参数")
)

// 服务端错误 500xx
var (
	ErrDB             = New(50001, "数据库错误")
	ErrCache          = New(50002, "缓存错误")
	ErrInternal       = New(50003, "内部错误")
)

// BizError 业务错误
type BizError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *BizError) Error() string {
	return e.Message
}

// New 创建业务错误
func New(code int, message string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
	}
}
