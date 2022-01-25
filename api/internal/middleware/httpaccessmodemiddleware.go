package middleware

import (
	"API_Gateway/api/internal/global"
	"API_Gateway/pkg/errcode"
	"context"
	"encoding/json"
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
		fmt.Println("请求成功")

		if serviceInfo, err := HTTPAccessMode(r, *global.SerInfo); err != nil {
			w.Write([]byte(err.Error()))
			fmt.Println("http proxy err ", err)
			return
		} else {
			// 如果能够从数据库获取这个请求对应的服务信息，就传递这个服务
			//fmt.Println("该请求匹配到的服务信息", Obj2Json(*serviceInfo))
			fmt.Println("匹配成功")
			ctx := context.WithValue(r.Context(), "serviceInfo", serviceInfo) // 这里key 最好使用自定义类型
			next(w, r.WithContext(ctx))                                       // WithContext返回r的浅层副本，其上下文更改为CTX
		}

	}
}

func Obj2Json(o interface{}) string {
	marshal, err := json.Marshal(o)
	if err != nil {
		fmt.Println("json marshal err", err)
	}
	return string(marshal)
}

// HTTPAccessMode 前端请求 与后端http服务   作用 : 找到请求对应的服务
func HTTPAccessMode(r *http.Request, info []global.ServiceDetail) (*global.ServiceDetail, error) {
	//fmt.Println("开始匹配接入方式")
	//fmt.Println("当前请求的path  HTTPAccessMode中间件", r.URL.Path)

	//1、前缀匹配 /abc ==> serviceSlice.rule
	//2、域名匹配 www.test.com ==> serviceSlice.rule
	//   根据请求可以得到该请求的的主机
	//   域名 host c.Request.Host
	//
	//   前缀 path c.Request.URL.Path
	//fmt.Println("r.RemoteAddr",r.RemoteAddr)
	//ip := r.Header.Get("X-Real-IP")
	//fmt.Println("获取到的ip:", ip)
	//r.Host = strings.Split(r.Host, ":")[0] //"www.baidu.com:8080" ++> "www.baidu.com" 去掉端口号

	//fmt.Println("当前请求去掉端口号的host", r.Host)
	//path := strings.TrimPrefix(r.URL.Path, "/")
	//host := r.Host
	//host = host[0:strings.Index(host, ":")]
	//fmt.Println("path", path)
	r.Header.Set("X-Real-IP", "127.0.0.1")
	// 1.找出所有的http服务
	for _, serviceItem := range info {
		if serviceItem.Info.LoadType != errcode.LoadTypeHTTP {
			continue
		}
		// 2.判断数据库这个服务是否是域名匹配
		if serviceItem.HTTPRule.RuleType == errcode.HTTPRuleTypeDomain {
			// 3.如果是域名匹配的话就判断是否和这个请求的host一样
			if serviceItem.HTTPRule.Rule == r.Host {
				return &serviceItem, nil
			}
		}
		// 3.判断这个服务如果是前缀匹配的话就判断是否和这个请求的path一样
		//fmt.Println("当前请求的ruleType", serviceItem.HTTPRule.RuleType)
		if serviceItem.HTTPRule.RuleType == errcode.HTTPRuleTypePrefixURL {
			// 判断当前这个请求是否是以这个前缀开头
			if strings.HasPrefix(r.URL.Path, serviceItem.HTTPRule.Rule) {
				//fmt.Println("数据库中对应的前缀", serviceItem.HTTPRule.Rule)
				//fmt.Println("r.当前请求的Path", r.URL.Path)
				return &serviceItem, nil
			}
		}
	}
	return nil, errors.New("该请求没有匹配到任何的服务")
}
