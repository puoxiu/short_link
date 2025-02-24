package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"short_link_pro/sl_convert/convert_api/internal/logic"
	"short_link_pro/sl_convert/convert_api/internal/svc"
	"short_link_pro/sl_convert/convert_api/internal/types"
)

func ConvertHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConvertRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewConvertLogic(r.Context(), svcCtx)
		resp, err := l.Convert(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
