package middleware

import (
	"API_Gateway/api/internal/global"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type UrlRewriteMiddleware struct {
}

func NewUrlRewriteMiddleware() *UrlRewriteMiddleware {
	return &UrlRewriteMiddleware{}
}

func (m *UrlRewriteMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := r.Context().Value("serviceInfo")
		serviceInfo := service.(*global.ServiceDetail)

		for _, item := range strings.Split(serviceInfo.HTTPRule.UrlRewrite, ",") {
			//fmt.Println("item rewrite",item)
			items := strings.Split(item, " ")
			if len(items) != 2 {
				continue
			}
			// 获取到正则表达式
			regexp, err := regexp.Compile(items[0])
			if err != nil {
				fmt.Println("regexp.Compile err", err)
				continue
			}

			fmt.Println("Url 重写之前的 path", r.URL.Path)
			fmt.Println("重写的规则", regexp)
			fmt.Println("转换后", items[1])
			replacePath := regexp.ReplaceAll([]byte(r.URL.Path), []byte(items[1]))
			r.URL.Path = string(replacePath)
			fmt.Println("Url 重写之后的 path", r.URL.Path)
		}
		next(w, r)
	}
}

// ^/http_proxy/url_rewrite(.*) /http_proxy/rewrite_url$1
// () 标记一个子表达式的开始和结束位置。子表达式可以获取供以后使用。要匹配这些字符
// .  匹配除换行符 \n 之外的任何单字符
// *  匹配前面的子表达式零次或多次
// $1 是第一个小括号里面的内容
// 只改变 url_rewrite
