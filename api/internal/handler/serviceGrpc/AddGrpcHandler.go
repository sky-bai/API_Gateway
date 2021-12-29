package serviceGrpc

import (
	"API_Gateway/util/reponse"
	"net/http"

	"API_Gateway/api/internal/logic/serviceGrpc"
	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func AddGrpcHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddGrpcRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := serviceGrpc.NewAddGrpcLogic(r.Context(), ctx)
		resp, err := l.AddGrpc(req)
		reponse.Response(w, resp, err) //②
	}
}
