package middleware

import (
	"API_Gateway/api/internal/global"
	"API_Gateway/pkg/errcode"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type HTTPAccessModeMiddleware struct {
}

func NewHTTPAccessModeMiddleware() *HTTPAccessModeMiddleware {
	return &HTTPAccessModeMiddleware{}
}

// 匹配接入方式 基于请求信息

// Handle 这个中间件去判断接入过来的请求是否后端有这个服务
func (m *HTTPAccessModeMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("哈哈哈👌", global.SerInfo)

		//fmt.Println("http proxy ", r)
		if serviceInfo, err := HTTPAccessMode(r, global.SerInfo); err != nil {
			fmt.Println("http proxy err ", err)
			return
		} else {
			ctx := context.WithValue(r.Context(), "serviceInfo", *serviceInfo)
			next(w, r.WithContext(ctx))
		}

	}
}

// HTTPAccessMode 前端请求 与后端http服务 想对接
func HTTPAccessMode(r *http.Request, info []global.ServiceDetail) (*global.ServiceDetail, error) {
	//1、前缀匹配 /abc ==> serviceSlice.rule
	//2、域名匹配 www.test.com ==> serviceSlice.rule
	//   根据请求可以得到该请求的的主机
	//   域名 host c.Request.Host
	//
	//   前缀 path c.Request.URL.Path

	r.Host = strings.Split(r.Host, ":")[0]
	fmt.Println("r.host", r.Host)
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/")
	//host := r.Host
	//host = host[0:strings.Index(host, ":")]

	path := r.URL.Path

	for _, serviceItem := range info {
		if serviceItem.Info.LoadType != errcode.LoadTypeHTTP {
			continue
		}
		if serviceItem.HTTPRule.RuleType == errcode.HTTPRuleTypeDomain {
			if serviceItem.HTTPRule.Rule == r.Host {
				fmt.Println("serviceItem.HTTPRule.Rule", serviceItem.HTTPRule.Rule)
				fmt.Println("r.Host", r.Host)
				return &serviceItem, nil
			}
		}
		if serviceItem.HTTPRule.RuleType == errcode.HTTPRuleTypePrefixURL {
			fmt.Println("serviceItem.HTTPRule.RuleType", serviceItem.HTTPRule.RuleType)
			if strings.HasPrefix(path, serviceItem.HTTPRule.Rule) {
				fmt.Println("serviceItem.HTTPRule.Rule", serviceItem.HTTPRule.Rule)
				fmt.Println("path", path)
				return &serviceItem, nil
			}
		}
	}
	return nil, errors.New("not matched service")
}
