package app

import (
	"API_Gateway/util/reponse"
	"net/http"

	"API_Gateway/api/internal/logic/app"
	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func AppDetailHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AppDetailRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := app.NewAppDetailLogic(r.Context(), ctx)
		resp, err := l.AppDetail(req)
		reponse.Response(w, resp, err) //②
	}
}
