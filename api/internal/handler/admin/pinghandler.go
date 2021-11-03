package admin

import (
	"API_Gateway/util/reponse"
	"net/http"

	"API_Gateway/api/internal/logic/admin"
	"API_Gateway/api/internal/svc"
)

func PingHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewPingLogic(r.Context(), ctx)
		resp, err := l.Ping()
		reponse.Response(w, resp, err)
	}
}
