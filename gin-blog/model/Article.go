package model

import (
	"gin-blog/utils/errmsg"
	"github.com/jinzhu/gorm"
)

type Article struct {
	Category Category `gorm:"foreignkey:Cid"`
	gorm.Model
	Title string `gorm:"type:varchar(100);not null" json:"title"`
	//文章所属类型的id
	Cid     int    `gorm:"type:int;not null" json:"cid"`
	Desc    string `gorm:"type:varchar(200)" json:"desc"`
	Content string `gorm:"type:longtext" json:"content"`
	Img     string `gorm:"type:varchar(100)" json:"img"`
}

// 查询某一分类下的所有文章
func GetCategoryArticle(id, pageSize, pageNum int) ([]Article, int, int) {
	var categoryarticle []Article
	var total int
	err := db.Preload("Category").Limit(pageSize).Offset(pageSize*(pageNum-1)).Where("cid = ?", id).Find(&categoryarticle).Count(&total).Error
	if err != nil {
		return nil, errmsg.ERROR_CATENAME_NOT_EXIST, 0
	}
	return categoryarticle, errmsg.SUCCESS, total
}

// 查询单个文章
func GetArticleInfo(id int) (Article, int) {
	var a Article
	err := db.Preload("Category").Where("id = ?", id).First(&a).Error
	if err != nil {
		return Article{}, errmsg.ERROR_ARTICLE_NOT_EXIST
	}
	return a, errmsg.SUCCESS

}

// 查询文章列表
func GetArticles(pageSize int, pageNum int) ([]Article, int, int) {
	var article []Article
	var total int
	//提前进行预加载，将关联到的表加载出来
	err = db.Preload("Category").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&article).Count(&total).Error
	if err != nil {
		return nil, errmsg.ERROR, 0
	}
	return article, errmsg.SUCCESS, total
}

// 新增文章
func CreateArticle(data *Article) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 编辑文章
func EditArticle(id int, data *Article) int {
	var article Article
	var maps = make(map[string]interface{})
	maps["title"] = data.Title
	maps["cid"] = data.Cid
	maps["desc"] = data.Desc
	maps["content"] = data.Content
	maps["img"] = data.Img
	err := db.Model(&article).Where("id = ?", id).Update(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 删除文章
func DeleteArticle(id int) int {
	var article Article
	err := db.Where("id =?", id).Delete(&article).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
