package handler

import (
	"net/http"

	"short_link_pro/sl_convert/convert_api/internal/logic"
	"short_link_pro/sl_convert/convert_api/internal/svc"
	"short_link_pro/sl_convert/convert_api/internal/types"

	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/zeromicro/go-zero/core/logx"
)

func ConvertHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConvertRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 参数校验--根据结构体validate tag
		if err := validator.New().StructCtx(r.Context(), &req); err != nil {
			logx.Errorf("validator check failed: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 执行业务逻辑	
		l := logic.NewConvertLogic(r.Context(), svcCtx)
		resp, err := l.Convert(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
