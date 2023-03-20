package v1

//控制api达到对数据库的操作，拿到model传回的code来对前端返回消息
import (
	"fmt"
	"gin-blog/model"
	"gin-blog/utils/errmsg"
	"gin-blog/utils/validater"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 标志码用于检验api的执行情况
var code int

// 添加用户
func AddUser(c *gin.Context) {
	//1.拿到前端的数据，入库进行检查看看是否存在
	var data model.User
	//一般来说绑定的错误不需要处理
	_ = c.ShouldBindJSON(&data)
	var msg string
	//进行数据验证
	msg, code = validater.Validate(data)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": msg,
		})
		return
	}
	//去检查用户是否存在
	code = model.CheckUser(data.UserName)

	if code == errmsg.SUCCESS {
		//无重名写入数据库
		code = model.CreateUser(&data)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

//查询单个用户

// 查询用户列表
func GetUsers(c *gin.Context) {
	//Atoi字符串转换成数字有两个返回值一个是数字另一个是err
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))

	if pageSize == 0 {
		//gorm提供方法如果给limit和offset传入-1就没有限制了
		pageSize = -1
	}
	if pageNum == 0 {
		pageNum = -1
	}
	data, total := model.GetUsers(pageSize, pageNum)
	code := errmsg.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data":   data,
		"msg":    errmsg.GetErrMsg(code),
		"total":  total,
	})
}

// 编辑用户
func EditUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var data model.User
	_ = c.ShouldBindJSON(&data)
	code = model.CheckUser(data.UserName)
	//用户名没有被使用
	if code == errmsg.SUCCESS {
		model.EditUser(id, &data)
	}
	//如果重名
	if code == errmsg.ERROR_USERNAME_USED {
		//调用 Abort 以确保这个请求的其他函数不会被调用。
		c.Abort()
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 删除用户
func DeleteUser(c *gin.Context) {
	//c.Param("id")返回的是一个字符串我们需要把它转换为int类型的数字
	id, _ := strconv.Atoi(c.Param("id"))
	code = model.DeleteUser(id)
	fmt.Println()
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
