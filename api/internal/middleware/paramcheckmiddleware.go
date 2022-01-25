package middleware

import (
	"net/http"
)

type ParamCheckMiddleware struct {
}

func NewParamCheckMiddleware() *ParamCheckMiddleware {
	return &ParamCheckMiddleware{}
}

//var Val *validator.Validate
//var Trans ut.Translator

func (m *ParamCheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//fmt.Println("执行参数校验中间件")
		//
		//zh := zh.New()
		//uni := ut.New(zh)
		//Trans, _ = uni.GetTranslator("zh")
		//Val = validator.New()
		//err := zh_translations.RegisterDefaultTranslations(Val, Trans)
		//if err != nil {
		//	fmt.Println("---", err)
		//	return
		//}
		//
		//Val.RegisterTranslation("required", Trans, func(ut ut.Translator) error {
		//	return ut.Add("required", "{0} 该字段不能为空!", true) // see universal-translator for details
		//}, func(ut ut.Translator, fe validator.FieldError) string {
		//	t, _ := ut.T("required", fe.Field())
		//	return t
		//})
		//
		//// 翻译器注册到validator
		//
		//Val.RegisterTagNameFunc(func(fld reflect.StructField) string {
		//	return fld.Tag.Get("comment")
		//})
		//
		//// 自定义验证方法
		////https://github.com/go-playground/validator/blob/v9/_examples/custom-validation/main.go
		//Val.RegisterValidation("valid_username", func(fl validator.FieldLevel) bool {
		//	return fl.Field().String() == "admin"
		//})
		//Val.RegisterValidation("valid_service_name", func(fl validator.FieldLevel) bool { // `^[a-zA-Z0-9_】$`
		//	matched, _ := regexp.Match(`^[a-zA-Z0-9_]{6,128}$`, []byte(fl.Field().String())) // 规定只能数字字母下划线 所以加^$
		//	return matched
		//})
		//Val.RegisterValidation("valid_rule", func(fl validator.FieldLevel) bool {
		//	matched, _ := regexp.Match(`^\S+$`, []byte(fl.Field().String())) // 非空 \S 非空白符 不能有空白符
		//	return matched
		//})
		//Val.RegisterValidation("valid_url_rewrite", func(fl validator.FieldLevel) bool { // 如果是空字符串就不进行校验
		//	if fl.Field().String() == "" {
		//		return true
		//	}
		//	for _, ms := range strings.Split(fl.Field().String(), ",") {
		//		if len(strings.Split(ms, " ")) != 2 {
		//			return false
		//		}
		//	}
		//	return true
		//})
		//Val.RegisterValidation("valid_header_transfor", func(fl validator.FieldLevel) bool {
		//	if fl.Field().String() == "" {
		//		return true
		//	}
		//	for _, ms := range strings.Split(fl.Field().String(), ",") {
		//		if len(strings.Split(ms, " ")) != 3 {
		//			return false
		//		}
		//	}
		//	return true
		//})
		//Val.RegisterValidation("valid_ipportlist", func(fl validator.FieldLevel) bool {
		//	for _, ms := range strings.Split(fl.Field().String(), ",") {
		//		if matched, _ := regexp.Match(`^\S+\:\d+$`, []byte(ms)); !matched {
		//			return false
		//		} // /d 0-9
		//	}
		//	return true
		//})
		//Val.RegisterValidation("valid_iplist", func(fl validator.FieldLevel) bool {
		//	if fl.Field().String() == "" {
		//		return true
		//	}
		//	for _, item := range strings.Split(fl.Field().String(), ",") {
		//		matched, _ := regexp.Match(`^\S+\:\d+$`, []byte(item)) //ip_addr  `^\S\:\d+$`
		//		if !matched {
		//			return false
		//		}
		//	}
		//	return true
		//})
		//Val.RegisterValidation("valid_weightlist", func(fl validator.FieldLevel) bool {
		//
		//	for _, ms := range strings.Split(fl.Field().String(), ",") {
		//		if matched, _ := regexp.Match(`^\d+$`, []byte(ms)); !matched {
		//			return false
		//		}
		//	}
		//	return true
		//})
		//
		////自定义翻译器
		////https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go
		//Val.RegisterTranslation("valid_username", Trans, func(ut ut.Translator) error {
		//	return ut.Add("valid_username", "{0} 填写不正确哦", true)
		//}, func(ut ut.Translator, fe validator.FieldError) string {
		//	t, _ := ut.T("valid_username", fe.Field())
		//	return t
		//})
		//Val.RegisterTranslation("valid_service_name", Trans, func(ut ut.Translator) error {
		//	return ut.Add("valid_service_name", "{0} 不符合输入格式", true)
		//}, func(ut ut.Translator, fe validator.FieldError) string {
		//	t, _ := ut.T("valid_service_name", fe.Field())
		//	return t
		//})
		//Val.RegisterTranslation("valid_rule", Trans, func(ut ut.Translator) error {
		//	return ut.Add("valid_rule", "{0} 需要填写", true)
		//}, func(ut ut.Translator, fe validator.FieldError) string {
		//	t, _ := ut.T("valid_rule", fe.Field())
		//	return t
		//})
		//Val.RegisterTranslation("valid_url_rewrite", Trans, func(ut ut.Translator) error {
		//	return ut.Add("valid_url_rewrite", "{0} 不符合输入格式", true)
		//}, func(ut ut.Translator, fe validator.FieldError) string {
		//	t, _ := ut.T("valid_url_rewrite", fe.Field())
		//	return t
		//})
		//Val.RegisterTranslation("valid_header_transfor", Trans, func(ut ut.Translator) error {
		//	return ut.Add("valid_header_transfor", "{0} 不符合输入格式", true)
		//}, func(ut ut.Translator, fe validator.FieldError) string {
		//	t, _ := ut.T("valid_header_transfor", fe.Field())
		//	return t
		//})
		//Val.RegisterTranslation("valid_ipportlist", Trans, func(ut ut.Translator) error {
		//	return ut.Add("valid_ipportlist", "{0} 不符合输入格式", true)
		//}, func(ut ut.Translator, fe validator.FieldError) string {
		//	t, _ := ut.T("valid_ipportlist", fe.Field())
		//	return t
		//})
		//Val.RegisterTranslation("valid_iplist", Trans, func(ut ut.Translator) error {
		//	return ut.Add("valid_iplist", "{0} 不符合输入格式", true)
		//}, func(ut ut.Translator, fe validator.FieldError) string {
		//	t, _ := ut.T("valid_iplist", fe.Field())
		//	return t
		//})
		//Val.RegisterTranslation("valid_weightlist", Trans, func(ut ut.Translator) error {
		//	return ut.Add("valid_weightlist", "{0} 不符合输入格式", true)
		//}, func(ut ut.Translator, fe validator.FieldError) string {
		//	t, _ := ut.T("valid_weightlist", fe.Field())
		//	return t
		//})

		next(w, r)
	}
}
