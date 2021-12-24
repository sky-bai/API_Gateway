package middleware

import (
	"API_Gateway/api/internal/global"
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

func (m *StripUrlMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := r.Context().Value("serviceInfo")
		serviceInfo := service.(*global.ServiceDetail)

		if serviceInfo.HTTPRule.RuleType == HTTPRuleTypePrefixURL || serviceInfo.HTTPRule.NeedStripUri == 1 {
			r.URL.Path = strings.Replace(r.URL.Path, serviceInfo.HTTPRule.Rule, "", 1)
		}

		next(w, r)
	}
}
