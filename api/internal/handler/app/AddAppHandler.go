package app

import (
	"API_Gateway/util/reponse"
	"net/http"

	"API_Gateway/api/internal/logic/app"
	"API_Gateway/api/internal/svc"
	"API_Gateway/api/internal/types"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func AddAppHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddAppRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := app.NewAddAppLogic(r.Context(), ctx)
		resp, err := l.AddApp(req)
		reponse.Response(w, resp, err) //②
	}
}
