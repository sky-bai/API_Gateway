package admin

import (
	"net/http"

	"API_Gateway/api/internal/logic/admin"
	"API_Gateway/api/internal/svc"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func PingHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewPingLogic(r.Context(), ctx)
		resp, err := l.Ping()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
