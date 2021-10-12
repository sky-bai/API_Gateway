package admin

import (
	"net/http"

	"API_Gateway/api/internal/logic/admin"
	"API_Gateway/api/internal/svc"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func AdminLogOutHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminLogOutLogic(r.Context(), ctx)
		resp, err := l.AdminLogOut()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
