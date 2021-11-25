package serviceInfo

import (
	"API_Gateway/util/reponse"
	"net/http"

	"API_Gateway/api/interna/logic/serviceInfo"
	"API_Gateway/api/interna/svc"
	"API_Gateway/api/interna/types"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func ServiceDetailHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ServiceDetailResquest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := serviceInfo.NewServiceDetailLogic(r.Context(), ctx)
		resp, err := l.ServiceDetail(req)
		reponse.Response(w, resp, err)
	}
}
