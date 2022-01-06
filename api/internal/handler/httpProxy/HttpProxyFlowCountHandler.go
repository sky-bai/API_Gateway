package httpProxy

import (
	"net/http"

	"API_Gateway/api/internal/logic/httpProxy"
	"API_Gateway/api/internal/svc"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func HttpProxyFlowCountHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := httpProxy.NewHttpProxyFlowCountLogic(r.Context(), ctx)
		resp, err := l.HttpProxyFlowCount()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
