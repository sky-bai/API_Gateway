package serviceInfo

import (
	"API_Gateway/util/reponse"
	"net/http"

	"API_Gateway/api/internal/logic/serviceInfo"
	"API_Gateway/api/internal/svc"
)

func ServiceFlowHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := serviceInfo.NewServiceFlowLogic(r.Context(), ctx)
		resp, err := l.ServiceFlow()
		reponse.Response(w, resp, err) //②
	}
}
