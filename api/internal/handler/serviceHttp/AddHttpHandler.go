package serviceHttp

import (
	"API_Gateway/util/reponse"
	"net/http"

	"API_Gateway/api/internal/logic/serviceHttp"
	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func AddHttpHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddHTTPResquest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := serviceHttp.NewAddHttpLogic(r.Context(), ctx)
		resp, err := l.AddHttp(req)
		reponse.Response(w, resp, err) //②
	}
}
