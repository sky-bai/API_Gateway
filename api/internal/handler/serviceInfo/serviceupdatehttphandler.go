package serviceInfo

import (
	"API_Gateway/util/reponse"
	"net/http"

	"API_Gateway/api/internal/logic/serviceInfo"
	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func ServiceUpdateHttpHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateHTTPResquest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := serviceInfo.NewServiceUpdateHttpLogic(r.Context(), ctx)
		resp, err := l.ServiceUpdateHttp(req)
		reponse.Response(w, resp, err)
	}
}
