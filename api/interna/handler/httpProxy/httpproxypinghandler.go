package httpProxy

import (
	"API_Gateway/util/reponse"
	"net/http"

	"API_Gateway/api/interna/logic/httpProxy"
	"API_Gateway/api/interna/svc"
)

func HttpProxyPingHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := httpProxy.NewHttpProxyPingLogic(r.Context(), ctx)
		resp, err := l.HttpProxyPing()
		reponse.Response(w, resp, err) //②
	}
}
