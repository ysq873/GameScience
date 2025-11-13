package handler

import (
	"net/http"
	"strconv"

	"backend/internal/repo"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getModelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
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
		repoM := repo.NewModelRepo(svcCtx.DB.Conn)
		m, err := repoM.GetByID(r.Context(), id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		cover := ""
		if m.CoverUrl.Valid {
			cover = m.CoverUrl.String
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
			"id":          m.Id,
			"title":       m.Title,
			"description": m.Description,
			"price_cents": m.PriceCents,
			"cover_url":   cover,
			"status":      m.Status,
		})
	}
}
