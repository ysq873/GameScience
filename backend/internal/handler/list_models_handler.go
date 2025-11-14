package handler

import (
    "net/http"
    "strconv"

    "backend/internal/middleware"
    "backend/internal/svc"
    "github.com/zeromicro/go-zero/rest/httpx"
    "backend/internal/repo"
)

type ModelItem struct {
    Id    int64  `json:"id"`
    Title string `json:"title"`
    Price int64  `json:"price_cents"`
    Cover string `json:"cover_url"`
    Status int `json:"status"`
}

func listModelsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        page, _ := strconv.Atoi(r.URL.Query().Get("page"))
        size, _ := strconv.Atoi(r.URL.Query().Get("size"))
        repoM := repo.NewModelRepo(svcCtx.DB.Conn)
        mine := r.URL.Query().Get("mine") == "1"
        var rows []repo.Model
        var err error
        if mine {
            sess, ok := middleware.GetSessionFromCtx(r.Context())
            if !ok {
                httpx.ErrorCtx(r.Context(), w, http.ErrNoCookie)
                return
            }
            rows, err = repoM.ListByOwner(r.Context(), sess.GetIdentity().Id, page, size)
        } else {
            rows, err = repoM.ListListed(r.Context(), page, size)
        }
        if err != nil { httpx.ErrorCtx(r.Context(), w, err); return }
        out := make([]ModelItem, 0, len(rows))
        for _, m := range rows {
            cover := ""
            if m.CoverUrl.Valid { cover = m.CoverUrl.String }
            out = append(out, ModelItem{Id: m.Id, Title: m.Title, Price: m.PriceCents, Cover: cover, Status: m.Status})
        }
        httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": out, "page": page, "size": size})
    }
}
