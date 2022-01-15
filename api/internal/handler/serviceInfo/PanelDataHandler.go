package serviceInfo

import (
	"API_Gateway/util/reponse"
	"net/http"

	"API_Gateway/api/internal/logic/serviceInfo"
	"API_Gateway/api/internal/svc"
)

func PanelDataHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := serviceInfo.NewPanelDataLogic(r.Context(), ctx)

		resp, err := l.PanelData()
		reponse.Response(w, resp, err) //②
	}
}
