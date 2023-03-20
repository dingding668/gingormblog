package v1

import (
	"gin-blog/middleware"
	"gin-blog/model"
	"gin-blog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	//对用户的权限进行验证
	var data model.User
	var token string
	_ = c.ShouldBindJSON(&data)
	code = model.CheckLogin(data.UserName, data.Password)
	if code == errmsg.SUCCESS {
		//成功登陆的话,设置token
		token, code = middleware.SetToken(data.UserName)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
		"token":   token,
	})
}
