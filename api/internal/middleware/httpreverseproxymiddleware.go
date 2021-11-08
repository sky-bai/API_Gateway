package middleware

import "net/http"

type HTTPReverseProxyMiddleware struct {
}

func NewHTTPReverseProxyMiddleware() *HTTPReverseProxyMiddleware {
	return &HTTPReverseProxyMiddleware{}
}

// 对接后端接口
func (m *HTTPReverseProxyMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Passthrough to next handler if need
		next(w, r)
	}
}
