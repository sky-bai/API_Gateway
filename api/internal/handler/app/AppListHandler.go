package app

import (
	"API_Gateway/util/reponse"
	"net/http"

	"API_Gateway/api/internal/logic/app"
	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func AppListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AppListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := app.NewAppListLogic(r.Context(), ctx)
		resp, err := l.AppList(req)
		reponse.Response(w, resp, err) //②
	}
}
