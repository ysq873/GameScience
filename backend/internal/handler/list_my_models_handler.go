package handler

import (
    "net/http"
    "strconv"

    "backend/internal/middleware"
    "backend/internal/repo"
    "backend/internal/svc"
    "github.com/zeromicro/go-zero/rest/httpx"
)

func listMyModelsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        page, _ := strconv.Atoi(r.URL.Query().Get("page"))
        size, _ := strconv.Atoi(r.URL.Query().Get("size"))
        sess, ok := middleware.GetSessionFromCtx(r.Context())
        if !ok {
            httpx.ErrorCtx(r.Context(), w, http.ErrNoCookie)
            return
        }
        repoM := repo.NewModelRepo(svcCtx.DB.Conn)
        rows, err := repoM.ListByOwner(r.Context(), sess.GetIdentity().Id, page, size)
        if err != nil { httpx.ErrorCtx(r.Context(), w, err); return }
        out := make([]map[string]interface{}, 0, len(rows))
        for _, m := range rows {
            cover := ""
            if m.CoverUrl.Valid { cover = m.CoverUrl.String }
            out = append(out, map[string]interface{}{
                "id": m.Id,
                "title": m.Title,
                "price_cents": m.PriceCents,
                "cover_url": cover,
                "status": m.Status,
            })
        }
        httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": out, "page": page, "size": size})
    }
}

