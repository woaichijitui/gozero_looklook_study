package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozero_looklook_study/app/usercenter/cmd/api/internal/logic/user"
	"gozero_looklook_study/app/usercenter/cmd/api/internal/svc"
	"gozero_looklook_study/app/usercenter/cmd/api/internal/types"
)

// register
func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}