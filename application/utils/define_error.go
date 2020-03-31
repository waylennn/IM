package utils


const (
	ErrCodeSuccess           = 0
	ErrCodeParameter         = 1001
	ErrCodeUserExist         = 1002
	ErrCodeServerBusy        = 1003
	ErrCodeUserNotExist      = 1004
	ErrCodeUserPasswordWrong = 1005
	ErrCodeCaptionHit        = 1006
	ErrCodeContentHit        = 1007
	ErrCodeNotLogin          = 1008
	ErrInvalidToken       = 1009
	ErrCodeGetTokenFailed       = 1011
	ErrCodeDb = 1012
	ErrRPC = 1013
	ErrCodeParamNotExist = 1014
	ErrConnTokenFailed = 1015

)

func GetMessage(code int) (message string) {
	switch code {
	case ErrCodeSuccess:
		message = "success"
	case ErrCodeParameter:
		message = "参数错误"
	case ErrCodeUserExist:
		message = "用户名已经存在"
	case ErrCodeServerBusy:
		message = "服务器繁忙"
	case ErrCodeUserNotExist:
		message = "用户名不存在"
	case ErrCodeUserPasswordWrong:
		message = "用户名或密码错误"
	case ErrCodeCaptionHit:
		message = "标题中含有非法内容, 请修改后发表"
	case ErrCodeContentHit:
		message = "内容中含有非法内容，请修改后发表"
	case ErrCodeNotLogin:
		message = "用户未登录"
	case ErrCodeGetTokenFailed:
		message = "生成token失败"
	case ErrInvalidToken:
		message = "token无效"
	case ErrCodeDb:
		message = "数据库相关错误"
	case ErrRPC:
		message = "RPC相关错误"
	case ErrCodeParamNotExist:
		message = "输入参数不完整"
	default:
		message = "未知错误"
	}
	return
}