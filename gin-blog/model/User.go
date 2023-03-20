package model

//对数据库的具体操作
import (
	"fmt"
	"gin-blog/utils/errmsg"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type User struct {
	gorm.Model //下面的json和前端传回的要一致
	//gorm的约定是蛇形的,数据库中UserName对应的是user_name
	//空字符串不是null所以也是可以传入数据库的
	//validate:"required"不能为空,字符串的长度大于等于4，小于等于12。
	UserName string `gorm:"type:varchar(20);not null" json:"username" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(100);not null" json:"password" validate:"required,min=6,max=20" label:"密码"`
	//1管理员，2阅读者    ;在添加用户的api中用户的角色吗必须大于等于2
	Role int `gorm:"type:int;DEFAULT:2" json:"role" validate:"required,gte=2" label:"角色码"`
}

// 查询用户是否存在
func CheckUser(name string) int {
	//实例化表
	var users User
	//在数据库中查询username等于用户注册的username如果存在就返回id，否则返回零值
	db.Select("id").Where("user_name = ?", name).First(&users)
	if users.ID > 0 {
		return errmsg.ERROR_USERNAME_USED //1001用户名重复
	}
	return errmsg.SUCCESS
}

// 新增用户,传入在控制器接收到的上下文，返回定义好的状态码
func CreateUser(data *User) int {
	//对密码进行加密
	//data.Password = ScryptPw(data.Password)
	err := db.Create(&data).Error
	if err != nil {
		//返回errmsg包里的信息
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 查询用户列表，涉及到分页
// 传入每页包含多少数据和需要的页数返回当前页的所有user数据
func GetUsers(pageSize int, pageNum int) ([]User, int) {
	var total int
	var users []User
	//db和err是db.go的全局变量，可以在本包访问
	//一页有几个然后需要偏移多少量
	err = db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Count(&total).Error
	if err != nil {
		//没有数据
		return nil, 0
	}
	return users, total
}

// 编辑用户,实际上就是更新数据
// 传入要修改的用户的id和要修改的数据(不包括密码)
func EditUser(id int, data *User) int {
	var user User
	var maps = make(map[string]interface{})
	maps["UserName"] = data.UserName
	maps["role"] = data.Role
	//Model 用来指定要操作的数据库模型
	err := db.Model(&user).Where("id = ?", id).Update(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 删除用户
func DeleteUser(id int) int {
	var user User
	err := db.Where("id =?", id).Delete(&user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 用户密码加密
// 采用钩子函数将密码保存,函数名是提前约定好的
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	//对密码进行加密
	u.Password = ScryptPw(u.Password)
	return nil
}

// 传入用户输入的密码，返回加密后的密码
func ScryptPw(password string) string {
	//加密的强度
	const cost = 10

	HashPw, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		//log.Fatal函数完成 1.打印输出内容 2.退出应用程序 3.defer函数不会执行
		log.Fatal(err)
	}
	return string(HashPw)
}

// 登录验证
func CheckLogin(username string, password string) int {
	var user User
	db.Where("user_name = ?", username).First(&user)
	fmt.Println(user)
	if user.ID == 0 {
		//用户不存在
		return errmsg.ERROR_USER_NOT_EXIST
	}
	PasswordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if PasswordErr != nil {
		//密码错误
		return errmsg.ERROR_PASSWORD_WRONG
	}
	if user.Role != 1 {
		//用户没有权限
		return errmsg.ERROR_USER_NO_ROIGHT
	}
	return errmsg.SUCCESS
}
