package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		//允许所有的域名访问
		AllowOrigins: []string{"*"},
		//允许的方法
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		//允许的header请求头字段
		AllowHeaders:  []string{"*", "Authorization"},
		ExposeHeaders: []string{"Content-Length", "text/plain", "Authorization", "Content-Type"},
		//预请求通过后12小时内不需要再次进行预请求
		MaxAge: 12 * time.Hour,
	},
	)
}
