package serviceGrpc

import (
	"net/http"

	"API_Gateway/api/interna/logic/serviceGrpc"
	"API_Gateway/api/interna/svc"
	"API_Gateway/api/interna/types"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func UpdateGrpcHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateGrpcRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := serviceGrpc.NewUpdateGrpcLogic(r.Context(), ctx)
		resp, err := l.UpdateGrpc(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
