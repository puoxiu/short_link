package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"short_link_pro/sl_auth/auth_api/internal/logic"
	"short_link_pro/sl_auth/auth_api/internal/svc"
	"short_link_pro/sl_auth/auth_api/internal/types"
)

func authenticationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AuthenticationRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewAuthenticationLogic(r.Context(), svcCtx)
		resp, err := l.Authentication(&req)
		// if err != nil {
		// 	httpx.ErrorCtx(r.Context(), w, err)
		// } else {
		// 	httpx.OkJsonCtx(r.Context(), w, resp)
		// }
		response(r, w, resp, err)
	}
}

func response(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	type Body struct {
		Code uint32 `json:"code"`
		Msg string `json:"msg"`
		Data interface{} `json:"data"`
	}
	if err == nil {
		// 成功
		b := &Body{
			Code: 0,
			Msg: "success",
			Data: resp,
		}
		httpx.WriteJson(w, http.StatusOK, b)
		return
	}
	// 失败
	errCode := uint32(10001)
	// errMsg := "苏武器error"
	httpx.WriteJson(w, http.StatusOK, &Body{
		Code: errCode,
		Msg: err.Error(),
		Data: nil,
	})
}
