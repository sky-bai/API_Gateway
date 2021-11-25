package admin

import (
	"net/http"

	"API_Gateway/api/interna/logic/admin"
	"API_Gateway/api/interna/svc"
	"API_Gateway/api/interna/types"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func AdminChangePwdHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FixPwdRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := admin.NewAdminChangePwdLogic(r.Context(), ctx)
		resp, err := l.AdminChangePwd(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
