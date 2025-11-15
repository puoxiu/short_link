package handler

import (
	"net/http"

	"short_link_pro/sl_redict/redict_api/internal/logic"
	"short_link_pro/sl_redict/redict_api/internal/svc"
	"short_link_pro/sl_redict/redict_api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)


func ShowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShowRequest
		err := httpx.Parse(r, &req) 
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewShowLogic(r.Context(), svcCtx)

		resp, err := l.Show(req.ShortUrl)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			// httpx.OkJsonCtx(r.Context(), w, resp)
			http.Redirect(w, r, resp.LongUrl, http.StatusFound)
		}
	}
}
