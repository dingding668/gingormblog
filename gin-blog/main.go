package main

import (
	"gin-blog/model"
	"gin-blog/routers"
)

func main() {
	//引用数据库
	model.InitDb()
	routers.InitRouter()
}
