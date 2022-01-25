package middleware

import (
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
	"reflect"
	"regexp"
	"strings"
)

type Validator struct {
	Validate  *validator.Validate
	Translate ut.Translator
}

var ValidatorHandler = NewValidator()

type ErrorString struct {
	errMessage string
}

func (e *ErrorString) Error() string {
	return e.errMessage
}

func NewValidator() *Validator {
	fmt.Println("NewValidator")
	zh := zh.New()
	uni := ut.New(zh)
	newValid := validator.New()
	newTrans, _ := uni.GetTranslator("zh")
	err := zh_translations.RegisterDefaultTranslations(newValid, newTrans)
	if err != nil {
		fmt.Println("---", err)
		return nil
	}

	// 如果字段为int类型，前端传的是0的话 这里required会有问题 它会认为0是没有填的
	newValid.RegisterTranslation("required", newTrans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} 该字段不能为空!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	// 翻译器注册到validator

	newValid.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("comment")
	})

	// 自定义验证方法
	//https://github.com/go-playground/validator/blob/v9/_examples/custom-validation/main.go
	err = newValid.RegisterValidation("valid_username", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "admin"
	})
	if err != nil {
		fmt.Println("---", err)
		return nil
	}
	newValid.RegisterValidation("valid_service_name", func(fl validator.FieldLevel) bool { // `^[a-zA-Z0-9_】$`
		// 判断是否满足正则表达式
		matched, _ := regexp.Match(`^[a-zA-Z0-9_]{6,128}$`, []byte(fl.Field().String())) // 规定只能数字字母下划线 所以加^$
		return matched
	})
	newValid.RegisterValidation("valid_rule", func(fl validator.FieldLevel) bool {
		matched, _ := regexp.Match(`^\S+$`, []byte(fl.Field().String())) // 非空 \S 非空白符 不能有空白符
		return matched
	})
	newValid.RegisterValidation("valid_url_rewrite", func(fl validator.FieldLevel) bool { // 如果是空字符串就不进行校验
		if fl.Field().String() == "" {
			return true
		}
		for _, ms := range strings.Split(fl.Field().String(), ",") {
			if len(strings.Split(ms, " ")) != 2 {
				return false
			}
		}
		return true
	})
	newValid.RegisterValidation("valid_header_transfor", func(fl validator.FieldLevel) bool {
		if fl.Field().String() == "" {
			return true
		}
		for _, ms := range strings.Split(fl.Field().String(), ",") {
			if len(strings.Split(ms, " ")) != 3 {
				return false
			}
		}
		return true
	})
	newValid.RegisterValidation("valid_ipportlist", func(fl validator.FieldLevel) bool {
		for _, ms := range strings.Split(fl.Field().String(), ",") {
			if matched, _ := regexp.Match(`^\S+\:\d+$`, []byte(ms)); !matched {
				return false
			} // /d 0-9
		}
		return true
	})
	newValid.RegisterValidation("valid_iplist", func(fl validator.FieldLevel) bool {
		if fl.Field().String() == "" {
			return true
		}
		for _, item := range strings.Split(fl.Field().String(), "\n") {
			matched, _ := regexp.Match(`^\S+$`, []byte(item)) //ip_addr  `^\S\:\d+$`  ^开始 $结束 \S非空白符 \d数字 +一个或多个
			if !matched {
				return false
			}
		}
		return true
	})
	newValid.RegisterValidation("valid_weightlist", func(fl validator.FieldLevel) bool {
		for _, ms := range strings.Split(fl.Field().String(), ",") {
			if matched, _ := regexp.Match(`^\d+$`, []byte(ms)); !matched {
				return false
			}
		}
		return true
	})

	//自定义翻译器
	//https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go
	newValid.RegisterTranslation("valid_username", newTrans, func(ut ut.Translator) error {
		return ut.Add("valid_username", "{0} 填写不正确哦", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_username", fe.Field())
		return t
	})
	newValid.RegisterTranslation("valid_service_name", newTrans, func(ut ut.Translator) error {
		return ut.Add("valid_service_name", "{0} 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_service_name", fe.Field())
		return t
	})
	newValid.RegisterTranslation("valid_rule", newTrans, func(ut ut.Translator) error {
		return ut.Add("valid_rule", "{0} 需要填写", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_rule", fe.Field())
		return t
	})
	newValid.RegisterTranslation("valid_url_rewrite", newTrans, func(ut ut.Translator) error {
		return ut.Add("valid_url_rewrite", "{0} 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_url_rewrite", fe.Field())
		return t
	})
	newValid.RegisterTranslation("valid_header_transfor", newTrans, func(ut ut.Translator) error {
		return ut.Add("valid_header_transfor", "{0} 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_header_transfor", fe.Field())
		return t
	})
	newValid.RegisterTranslation("valid_ipportlist", newTrans, func(ut ut.Translator) error {
		return ut.Add("valid_ipportlist", "{0} 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_ipportlist", fe.Field())
		return t
	})
	newValid.RegisterTranslation("valid_iplist", newTrans, func(ut ut.Translator) error {
		return ut.Add("valid_iplist", "{0} 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_iplist", fe.Field())
		return t
	})
	newValid.RegisterTranslation("valid_weightlist", newTrans, func(ut ut.Translator) error {
		return ut.Add("valid_weightlist", "权重列表 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_weightlist", fe.Field())
		return t
	})
	return &Validator{
		Validate:  newValid,
		Translate: newTrans,
	}
}
