package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-cloud-disk/core/internal/logic"
	"go-cloud-disk/core/internal/svc"
	"go-cloud-disk/core/internal/types"
)

func FileUploadChunkCompleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileUploadChunkCompleteRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		l := logic.NewFileUploadChunkCompleteLogic(r.Context(), svcCtx)
		resp, err := l.FileUploadChunkComplete(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
