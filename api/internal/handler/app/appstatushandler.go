package app

import (
	"net/http"

	"API_Gateway/api/internal/logic/app"
	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func AppStatusHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AppStatusRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := app.NewAppStatusLogic(r.Context(), ctx)
		resp, err := l.AppStatus(req)
		reponse.Response(w, resp, err) //②
	}
}
