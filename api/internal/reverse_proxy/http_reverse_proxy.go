package reverse_proxy

import (
	"API_Gateway/api/internal/reverse_proxy/load_balance"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewLoadBalanceReverseProxy(c *http.Request, lb load_balance.LoadBalance, trans *http.Transport) *httputil.ReverseProxy {
	//请求协调者
	director := func(req *http.Request) {
		// 1.更新下游服务器
		lb.Update()
		// 2.获取下游服务器
		nextAddr, err := lb.Get(req.URL.String())
		fmt.Println(req.URL.String())
		logx.Error("下游服务器地址", nextAddr)
		//fmt.Println("当前请求的下游服务器地址",nextAddr)
		//todo 优化点3
		if err != nil || nextAddr == "" {
			panic("get next addr fail")
		}
		target, err := url.Parse(nextAddr)
		if err != nil {
			panic(err)
		}

		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		//fmt.Println("请求的URL",req.URL.Path)
		//fmt.Println("目标的path",target.Path)

		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path) // 在target.Path后面拼接req.URL.Path
		//fmt.Println("代理中拼接的URL",req.URL.Path)
		req.Host = target.Host
		//fmt.Println("代理中获取的Scheme",req.URL.Scheme)
		//fmt.Println("代理中获取的Host",req.URL.Host)
		//fmt.Println("代理中的Path",req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}
		fmt.Println("代理请求实际服务器的URL", req.URL.String())
	}

	//更改内容
	modifyFunc := func(resp *http.Response) error {
		//if strings.Contains(resp.Header.Get("Connection"), "Upgrade") {
		//	return nil
		//}

		//todo 优化点2
		//var payload []byte
		//var readErr error
		//
		//if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		//	gr, err := gzip.NewReader(resp.Body)
		//	if err != nil {
		//		return err
		//	}
		//	payload, readErr = ioutil.ReadAll(gr)
		//	resp.Header.Del("Content-Encoding")
		//} else {
		//	payload, readErr = ioutil.ReadAll(resp.Body)
		//}
		//if readErr != nil {
		//	return readErr
		//}
		//
		//c.Set("status_code", resp.StatusCode)
		//c.Set("payload", payload)
		//resp.Body = ioutil.NopCloser(bytes.NewBuffer(payload))
		//resp.ContentLength = int64(len(payload))
		//resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(payload)), 10))

		//payload := "hello world"
		//resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(payload)))
		//resp.ContentLength = int64(len(payload))
		//resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(payload)), 10))

		return nil
	}

	//错误回调 ：关闭real_server时测试，错误回调
	//范围：transport.RoundTrip发生的错误、以及ModifyResponse发生的错误
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		w.Write([]byte("请求失败，错误信息为 " + err.Error()))
	}
	return &httputil.ReverseProxy{Director: director, ModifyResponse: modifyFunc, Transport: trans, ErrorHandler: errFunc}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
