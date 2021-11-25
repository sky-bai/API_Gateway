package admin

import (
	"net/http"

	"API_Gateway/api/interna/logic/admin"
	"API_Gateway/api/interna/svc"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func AdminInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminInfoLogic(r.Context(), ctx)
		resp, err := l.AdminInfo()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
