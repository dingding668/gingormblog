package middleware

import (
	"gin-blog/utils"
	"gin-blog/utils/errmsg"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// 密钥参数
var JwtKey = []byte(utils.JwtKey)

// 额外要记录的信息
type MyClaims struct {
	//必须和user模型里面的用户名保持一致
	UserName string `json:"username"`
	jwt.StandardClaims
}

// 生成token
func SetToken(username string) (string, int) {
	//生成时间
	expireTime := time.Now().Add(10 * time.Hour)
	Setclaims := MyClaims{
		username,
		jwt.StandardClaims{
			//有效时间是十个小时，Unix()是一个时间戳，是自1970 年 1 月 1 日（00:00:00 GMT）以来的秒数
			ExpiresAt: expireTime.Unix(),
			//签发人
			Issuer: "dingding",
		},
	}

	//两个参数一个是签发的方法一个是设计好的参数，会返回一个token的指针
	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, Setclaims)

	token, err := reqClaim.SignedString(JwtKey)
	if err != nil {
		return "", errmsg.ERROR
	}
	return token, errmsg.SUCCESS
}

// 验证token
func CheckToken(token string) (*MyClaims, int) {
	setToken, _ := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if key, _ := setToken.Claims.(*MyClaims); setToken.Valid { //settoken.Valid token是否有效
		//如果成功就返回，我们要记录的信息
		return key, errmsg.SUCCESS
	} else {
		return nil, errmsg.ERROR
	}

}

var code int

// jwt中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		tokenHeader := c.Request.Header.Get("Authorization")

		if tokenHeader == "" {
			code = errmsg.ERROR_TOKEN_EXIST
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		checkToken := strings.SplitN(tokenHeader, " ", 2)
		if len(checkToken) != 2 && checkToken[0] != "Bearer" {
			code := errmsg.ERROR_TOKEN_TYPEWRONG
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		//checkToken中的第二个元素进行验证
		key, tcode := CheckToken(checkToken[1])

		if tcode == errmsg.ERROR {
			code = errmsg.ERROR_TOKEN_WRONG
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		//现在的时间大于有效时间
		if time.Now().Unix() > key.ExpiresAt {
			code = errmsg.ERROR_TOKEN_RUNTIME
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		//验证都通过了，就把用户名绑定到username上去
		c.Set("username", key.UserName)
		//进入下一步
		c.Next()
	}
}
