package serviceTcp

import (
	"API_Gateway/util/reponse"
	"net/http"

	"API_Gateway/api/interna/logic/serviceTcp"
	"API_Gateway/api/interna/svc"
	"API_Gateway/api/interna/types"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func UpdateTcpHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateTcpRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := serviceTcp.NewUpdateTcpLogic(r.Context(), ctx)
		resp, err := l.UpdateTcp(req)
		reponse.Response(w, resp, err) //②
	}
}
