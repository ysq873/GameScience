package handler

import (
    "encoding/json"
    "net/http"
    "strconv"

    "backend/internal/middleware"
    "backend/internal/repo"
    "backend/internal/svc"

    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/rest/httpx"
)

func addFavoriteModelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        sess, ok := middleware.GetSessionFromCtx(r.Context())
        if !ok { http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized); return }
        uid := sess.GetIdentity().Id
        midStr := r.URL.Query().Get("model_id")
        if midStr == "" { _ = r.ParseForm(); midStr = r.FormValue("model_id") }
        if midStr == "" {
            var m map[string]interface{}
            _ = json.NewDecoder(r.Body).Decode(&m)
            if v, ok := m["model_id"]; ok { midStr = strconv.FormatInt(int64(getNumber(v)), 10) }
        }
        if midStr == "" { http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest); return }
        mid, err := strconv.ParseInt(midStr, 10, 64)
        if err != nil { http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest); return }
        fr := repo.NewFavoriteRepo(svcCtx.DB.Conn)
        logx.Infof("favorite add uid=%s model_id=%d", uid, mid)
        if err := fr.Add(r.Context(), uid, mid); err != nil { httpx.ErrorCtx(r.Context(), w, err); return }
        httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{ "ok": true })
    }
}

func removeFavoriteModelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        sess, ok := middleware.GetSessionFromCtx(r.Context())
        if !ok { http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized); return }
        uid := sess.GetIdentity().Id
        midStr := r.URL.Query().Get("model_id")
        if midStr == "" { _ = r.ParseForm(); midStr = r.FormValue("model_id") }
        if midStr == "" {
            var m map[string]interface{}
            _ = json.NewDecoder(r.Body).Decode(&m)
            if v, ok := m["model_id"]; ok { midStr = strconv.FormatInt(int64(getNumber(v)), 10) }
        }
        if midStr == "" { http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest); return }
        mid, err := strconv.ParseInt(midStr, 10, 64)
        if err != nil { http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest); return }
        fr := repo.NewFavoriteRepo(svcCtx.DB.Conn)
        logx.Infof("favorite remove uid=%s model_id=%d", uid, mid)
        if err := fr.Remove(r.Context(), uid, mid); err != nil { httpx.ErrorCtx(r.Context(), w, err); return }
        httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{ "ok": true })
    }
}

func listFavoriteModelsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        sess, ok := middleware.GetSessionFromCtx(r.Context())
        if !ok { http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized); return }
        uid := sess.GetIdentity().Id
        fr := repo.NewFavoriteRepo(svcCtx.DB.Conn)
        rows, err := fr.ListModelInfosByUser(r.Context(), uid)
        if err != nil { httpx.ErrorCtx(r.Context(), w, err); return }
        out := make([]map[string]interface{}, 0, len(rows))
        for _, m := range rows {
            cover := ""
            if m.Cover.Valid { cover = m.Cover.String }
            out = append(out, map[string]interface{}{
                "id": m.Id,
                "title": m.Title,
                "price_cents": m.Price,
                "cover_url": cover,
            })
        }
        httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": out})
    }
}

func getNumber(v interface{}) int {
    switch t := v.(type) {
    case float64:
        return int(t)
    case int:
        return t
    case string:
        if n, err := strconv.Atoi(t); err == nil { return n }
    }
    return 0
}