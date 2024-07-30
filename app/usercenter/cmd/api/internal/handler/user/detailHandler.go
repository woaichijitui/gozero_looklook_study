package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozero_looklook_study/app/usercenter/cmd/api/internal/logic/user"
	"gozero_looklook_study/app/usercenter/cmd/api/internal/svc"
	"gozero_looklook_study/app/usercenter/cmd/api/internal/types"
)

// get user info
func DetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewDetailLogic(r.Context(), svcCtx)
		resp, err := l.Detail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
