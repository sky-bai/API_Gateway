package httpsProxy

import (
	"API_Gateway/util/reponse"
	"net/http"

	"API_Gateway/api/internal/logic/httpsProxy"
	"API_Gateway/api/internal/svc"
)

func HttpsProxyPingHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := httpsProxy.NewHttpsProxyPingLogic(r.Context(), ctx)
		resp, err := l.HttpsProxyPing()
		reponse.Response(w, resp, err)
	}
}
