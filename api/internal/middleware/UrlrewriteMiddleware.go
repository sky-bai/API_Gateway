package middleware

import (
	"API_Gateway/api/internal/global"
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
			regexp, err := regexp.Compile(items[0])
			if err != nil {
				//fmt.Println("regexp.Compile err",err)
				continue
			}
			//fmt.Println("before rewrite",c.Request.URL.Path)
			replacePath := regexp.ReplaceAll([]byte(r.URL.Path), []byte(items[1]))
			r.URL.Path = string(replacePath)
			//fmt.Println("after rewrite",c.Request.URL.Path)
		}
		next(w, r)
	}
}