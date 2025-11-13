package handler

import (
    "net/http"
    "strconv"

    "backend/internal/svc"
    "github.com/zeromicro/go-zero/rest/httpx"
    "backend/internal/repo"
)

type ModelItem struct {
    Id    int64  `json:"id"`
    Title string `json:"title"`
    Price int64  `json:"price_cents"`
    Cover string `json:"cover_url"`
}

func listModelsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        page, _ := strconv.Atoi(r.URL.Query().Get("page"))
        size, _ := strconv.Atoi(r.URL.Query().Get("size"))
        repoM := repo.NewModelRepo(svcCtx.DB.Conn)
        rows, err := repoM.ListListed(r.Context(), page, size)
        if err != nil { httpx.ErrorCtx(r.Context(), w, err); return }
        out := make([]ModelItem, 0, len(rows))
        for _, m := range rows {
            cover := ""
            if m.CoverUrl.Valid { cover = m.CoverUrl.String }
            out = append(out, ModelItem{Id: m.Id, Title: m.Title, Price: m.PriceCents, Cover: cover})
        }
        httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": out, "page": page, "size": size})
    }
}
