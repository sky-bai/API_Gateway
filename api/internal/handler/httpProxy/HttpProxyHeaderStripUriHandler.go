package httpProxy

import (
	"net/http"

	"API_Gateway/api/internal/logic/httpProxy"
	"API_Gateway/api/internal/svc"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func HttpProxyHeaderStripUriHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := httpProxy.NewHttpProxyHeaderStripUriLogic(r.Context(), ctx)
		resp, err := l.HttpProxyHeaderStripUri()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
