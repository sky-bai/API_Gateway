package middleware

import (
	"API_Gateway/api/internal/global"
	"fmt"
	"net/http"
	"strings"
)

type StripUrlMiddleware struct {
}

func NewStripUrlMiddleware() *StripUrlMiddleware {
	return &StripUrlMiddleware{}
}

const (
	HTTPRuleTypePrefixURL = 0
	HTTPRuleTypeDomain    = 1
)

// Handle 如果请求是前缀接入并且需要去除数据库中配置的路径的话 就去除数据库中配置的路径
func (m *StripUrlMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1.获取本次请求对应的配置
		service := r.Context().Value("serviceInfo")
		serviceInfo := service.(*global.ServiceDetail)

		// 2.如果是前缀接入的话并且需要去除后缀的话 就去除后缀
		if serviceInfo.HTTPRule.RuleType == HTTPRuleTypePrefixURL && serviceInfo.HTTPRule.NeedStripUri == 1 {
			fmt.Println("请求前的路径：", r.URL.Path)
			fmt.Println("数据库中对应的路径", serviceInfo.HTTPRule.Rule)
			r.URL.Path = strings.Replace(r.URL.Path, serviceInfo.HTTPRule.Rule, "", 1)
			fmt.Println("修改后的路径：", r.URL.Path)
		}
		next(w, r)
	}
}
