package model

import (
	"fmt"
	"gin-blog/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

//数据库的入口文件，用来建立数据连接的参数

var db *gorm.DB
var err error

func InitDb() {
	path := strings.Join([]string{utils.Dbuser, ":", utils.DbPassword, "@tcp(", utils.DbHost, utils.DbPost, ")/", utils.DbName, "?charset=utf8mb4&parseTime=True&loc=Local"}, "")
	db, err = gorm.Open(utils.Db, path)
	if err != nil {
		fmt.Println("连接数据库失败", err)
	}
	//禁用默认表名的复数形式 结构体user 数据库users
	db.SingularTable(true)
	//自动迁移只会创建表、缺失的列、缺失的索引，不会改变现有列的类型或者删除未使用的列
	//迁移自己的模型
	db.AutoMigrate(&User{}, &Article{}, &Category{})
	//参数设置
	//设置连接池中的最大闲置连接数
	db.DB().SetMaxIdleConns(10)
	//设置连接池中的最大连接数量
	db.DB().SetMaxOpenConns(100)
	//设置连接的最大可复用时间，建议不要超过gin连接的最长时间,10秒
	db.DB().SetConnMaxLifetime(10 * time.Second)

}
