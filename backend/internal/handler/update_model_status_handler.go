package handler

import (
	"net/http"
	"strconv"

	"backend/internal/repo"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type UpdateStatusReq struct {
	Status string `json:"status"`
}

func updateModelStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		var req UpdateStatusReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if req.Status != "listed" && req.Status != "delisted" && req.Status != "draft" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		repoM := repo.NewModelRepo(svcCtx.DB.Conn)
		if err := repoM.UpdateStatus(r.Context(), id, req.Status); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"ok": true})
	}
}
