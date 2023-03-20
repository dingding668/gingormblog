package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"time"
)

func Logger() gin.HandlerFunc {
	//日志存放路径
	filePath := "log/log"
	src, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("err:", err)
	}
	logger := logrus.New()
	//输出到指定路径
	logger.Out = src
	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	//日志分割
	logWriter, _ := rotatelogs.New(
		//日志文件名称
		filePath+"%Y%m%d.log",
		//最大保存时间，一周
		rotatelogs.WithMaxAge(7*24*time.Hour),
		//什么时候分割一次
		rotatelogs.WithRotationTime(24*time.Hour),
		//retalog.WithLinkName(linkName),
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	//按时间分割日志
	logrus.AddHook(hook)

	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		stopTime := time.Since(startTime)
		spendTime := fmt.Sprintf("%d ms", int(math.Ceil(float64(stopTime.Nanoseconds())/1000000.0)))
		hostName, err := os.Hostname()
		if err != nil {
			//未知主机访问
			hostName = "unknown"
		}
		//状态码
		statusCode := c.Writer.Status()
		//客户端ip
		clientIp := c.ClientIP()
		//使用什么浏览器
		userAgent := c.Request.UserAgent()
		//请求大小
		dataSize := c.Writer.Size()
		//请求的方法
		method := c.Request.Method
		//请求路径
		path := c.Request.RequestURI

		entry := logger.WithFields(logrus.Fields{
			"HostName":  hostName,
			"status":    statusCode,
			"SpendTime": spendTime,
			"IP":        clientIp,
			"Method":    method,
			"Path":      path,
			"DataSize":  dataSize,
			"Agent":     userAgent,
		})
		//记录系统内部有错误
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode >= 400 {
			entry.Warn()
		} else {
			entry.Info()
		}

	}
}
