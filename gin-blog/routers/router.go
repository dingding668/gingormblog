package routers

import (
	v1 "gin-blog/api/v1"
	"gin-blog/middleware"
	"gin-blog/utils"
	"github.com/gin-gonic/gin"
)

// 入口路由
func InitRouter() {
	//设置qppmode
	gin.SetMode(utils.AppMode)
	//初始化路由
	//这里用default和new都可以区别就是default增加了两个中间件
	r := gin.New()
	r.Use(gin.Recovery())
	//使用日志中间件
	r.Use(middleware.Logger())
	//解决跨域问题
	r.Use(middleware.Cors())
	//Auth表示需要鉴权的
	auth := r.Group("api/v1")
	auth.Use(middleware.JwtToken())
	{
		//用户模块的路由接口
		//查询用户列表
		auth.GET("users", v1.GetUsers)
		//编辑用户,参数携带在路径中
		auth.POST("user/:id", v1.EditUser)
		//删除用户
		auth.GET("user/:id", v1.DeleteUser)

		//分类模块的路由接口
		//添加分类
		auth.POST("category/add", v1.AddCategory)

		//编辑分类
		auth.POST("category/:id", v1.EditCategory)
		//删除分类
		auth.GET("category/:id", v1.DeleteCategory)

		//文章模块的路由接口
		//添加文章
		auth.POST("article/add", v1.AddArticle)
		//编辑文章
		auth.POST("article/:id", v1.EditArticle)
		//删除文章
		auth.GET("article/:id", v1.DeleteArticle)
		//上传文件
		//auth.POST("upload", v1.Upload)
	}
	//公共接口
	router := r.Group("api/v1")
	{
		//添加用户
		router.POST("user/add", v1.AddUser)
		//查询分类
		router.GET("category", v1.GetCategorys)
		//查询文章列表
		router.GET("article", v1.GetArticles)
		//查询某一分类下的所有文章
		router.GET("article/list/:id", v1.GetCategoryArticle)
		//查询某一文章的具体信息
		router.GET("article/info/:id", v1.GetArticleInfo)
		//用户登录
		router.POST("login", v1.Login)
	}

	//有两种方式跑起来一种是返回gin引擎，然后再main函数中调用
	//也可以直接run就会跑起来
	_ = r.Run(utils.HttpPort)
}
