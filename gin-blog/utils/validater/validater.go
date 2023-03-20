package validater

import (
	"fmt"
	"gin-blog/utils/errmsg"
	"github.com/go-playground/locales/zh_Hans_CN"
	unTrans "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

//如果验证不通过，默认是英文，可以翻译成不同的语言

//不知道传入的参数是什么类型就可以用接口然后断言
func Validate(data interface{}) (string, int) {
	validate := validator.New()

	uni := unTrans.New(zh_Hans_CN.New())
	//翻译方法
	trans, _ := uni.GetTranslator("zh_Hans_CN")
	//注册翻译方法
	err := zhTrans.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		fmt.Println("err:", err)
	}

	//将label标签进行映射
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		return label
	})

	//对传入的data进行验证
	err = validate.Struct(data)
	if err != nil {
		//将错误循环出来
		for _, v := range err.(validator.ValidationErrors) {
			return v.Translate(trans), errmsg.ERROR
		}
	}
	return "", errmsg.SUCCESS
}
