package middleware

import (
	"API_Gateway/api/internal/global"
	"fmt"
	"net/http"
	"strings"
)

type HeaderTransferMiddleware struct {
}

func NewHeaderTransferMiddleware() *HeaderTransferMiddleware {
	return &HeaderTransferMiddleware{}
}

// Handle 修改请求重写层的头部信息
func (m *HeaderTransferMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 1.获取该请求的详细信息
		service := r.Context().Value("serviceInfo")
		serviceInfo := service.(*global.ServiceDetail)

		// 2.header Transform 是多个header ( add HeaderName HeaderValue,add HeaderName HeaderValue )
		for _, value := range strings.Split(serviceInfo.HTTPRule.HeaderTransfor, ",") { // 以逗号分割出 多个操作header头的信息
			item := strings.Split(value, " ") // 以空格分割出 headerName 和 headerValue
			if len(item) != 3 {               // 验证该信息是否有效
				continue
			}
			// 3.添加header
			if item[0] == "add" || item[0] == "edit" {
				r.Header.Set(item[1], item[2])
				fmt.Println("添加成功", item[1], item[2])
			}
			// 4.删除header
			if item[0] == "del" {
				r.Header.Del(item[1])
				fmt.Println("删除成功", item[1])
			}

		}
		next(w, r)
	}
}
