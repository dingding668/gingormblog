package errmsg

// 错误处理的包
// 常量
const (
	//成功
	SUCCESS = 200
	ERROR   = 500

	//code = 1000... 表示用户模块的错误
	ERROR_USERNAME_USED   = 1001 //重名
	ERROR_PASSWORD_WRONG  = 1002 //密码错误
	ERROR_USER_NOT_EXIST  = 1003 //用户不存在
	ERROR_TOKEN_EXIST     = 1004 //用户携带token不存在
	ERROR_TOKEN_RUNTIME   = 1005 //用户token超时
	ERROR_TOKEN_WRONG     = 1006 //用户携带的token错误
	ERROR_TOKEN_TYPEWRONG = 1007 //token格式错误
	ERROR_USER_NO_ROIGHT  = 1008 //用户没有权限
	//code = 2000... 表示分类模块的错误
	ERROR_CATENAME_USED      = 2001
	ERROR_CATENAME_NOT_EXIST = 2002
	//code = 3000... 表示文章模块的错误
	ERROR_ARTICLE_NOT_EXIST = 3001
)

// 错误的编码和信息的映射
var codeMsg = map[int]string{
	SUCCESS:                  "OK",
	ERROR:                    "FAIL",
	ERROR_USERNAME_USED:      "用户名已经存在！",
	ERROR_PASSWORD_WRONG:     "密码错误！",
	ERROR_USER_NOT_EXIST:     "用户不存在！",
	ERROR_TOKEN_EXIST:        "TOKEN不存在！",
	ERROR_TOKEN_RUNTIME:      "TOKEN已经过期！",
	ERROR_TOKEN_WRONG:        "TOKEN不正确！",
	ERROR_TOKEN_TYPEWRONG:    "TOKEN格式错误",
	ERROR_CATENAME_USED:      "该分类已经存在",
	ERROR_ARTICLE_NOT_EXIST:  "文章不存在",
	ERROR_CATENAME_NOT_EXIST: "该分类不存在",
	ERROR_USER_NO_ROIGHT:     "该用户没有权限",
}

func GetErrMsg(code int) string {
	//根据错误码返回错误
	return codeMsg[code]
}
