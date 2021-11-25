package admin

import (
	"API_Gateway/util/reponse"
	"net/http"

	"API_Gateway/api/interna/logic/admin"
	"API_Gateway/api/interna/svc"
	"API_Gateway/api/interna/types"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func AdminLoginHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := admin.NewAdminLoginLogic(r.Context(), ctx)
		resp, err := l.AdminLogin(req)
		reponse.Response(w, resp, err)
	}
}
