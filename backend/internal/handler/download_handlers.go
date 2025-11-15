package handler

import (
    "database/sql"
    "errors"
    "encoding/json"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "time"
    "strings"

    "backend/internal/middleware"
    "backend/internal/repo"
    "backend/internal/svc"

    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/zeromicro/go-zero/rest/httpx"
    "github.com/zeromicro/go-zero/core/logx"
)

func listPurchasesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, ok := middleware.GetSessionFromCtx(r.Context())
        if !ok {
            http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
            return
        }
		uid := sess.GetIdentity().Id
		type Row struct {
			ModelId  int64
			Title    sql.NullString
			CoverUrl sql.NullString
		}
		var rows []Row
		err := svcCtx.DB.Conn.QueryRowsCtx(r.Context(), &rows, "SELECT p.model_id, m.title, m.cover_url FROM purchases p JOIN models m ON m.id=p.model_id WHERE p.user_id=? ORDER BY p.id DESC", uid)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		out := make([]map[string]interface{}, 0, len(rows))
		for _, r := range rows {
			title := ""
			if r.Title.Valid {
				title = r.Title.String
			}
			cover := ""
			if r.CoverUrl.Valid {
				cover = r.CoverUrl.String
			}
			out = append(out, map[string]interface{}{"model_id": r.ModelId, "title": title, "cover_url": cover})
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": out})
	}
}

func generateDownloadTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, ok := middleware.GetSessionFromCtx(r.Context())
        if !ok {
            http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
            return
        }
		uid := sess.GetIdentity().Id
		midStr := r.URL.Query().Get("model_id")
		if midStr == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		mid, err := strconv.ParseInt(midStr, 10, 64)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// ensure purchased
		var cnt int
		_ = svcCtx.DB.Conn.QueryRowCtx(r.Context(), &cnt, "SELECT COUNT(1) FROM purchases WHERE user_id=? AND model_id=?", uid, mid)
		if cnt == 0 {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		// ensure file exists before issuing token
		var rel string
		_ = svcCtx.DB.Conn.QueryRowCtx(r.Context(), &rel, "SELECT file_path FROM models WHERE id=?", mid)
		if rel == "" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		base := filepath.Clean(svcCtx.Config.StorageBasePath)
		abs := filepath.Clean(filepath.Join(base, filepath.FromSlash(rel)))
		if !strings.HasPrefix(abs, base) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		if fi, statErr := os.Stat(abs); statErr != nil || fi.IsDir() {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		dr := repo.NewDownloadRepo(svcCtx.DB.Conn)
		token, err := dr.CreateToken(r.Context(), uid, mid, 10*time.Minute)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"token": token})
	}
}

func downloadByTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        sess, ok := middleware.GetSessionFromCtx(r.Context())
    if !ok {
        http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
        return
    }
        uid := sess.GetIdentity().Id
        token := r.URL.Query().Get("token")
        if token == "" {
            http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
            return
        }
        dr := repo.NewDownloadRepo(svcCtx.DB.Conn)
        mid, err := dr.ConsumeToken(r.Context(), token, uid)
        if err != nil {
            if errors.Is(err, sqlx.ErrNotFound) || errors.Is(err, sql.ErrNoRows) || errors.Is(err, repo.ErrTokenExpired) {
                logx.Errorf("download not found: uid=%s token=%s err=%v", uid, token, err)
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                code := "token_not_found"
                if errors.Is(err, repo.ErrTokenExpired) { code = "token_expired" }
                _ = json.NewEncoder(w).Encode(map[string]interface{}{"code": code, "message": http.StatusText(http.StatusNotFound)})
                return
            }
            if errors.Is(err, repo.ErrTokenUserMismatch) {
                logx.Errorf("download forbidden: uid=%s token=%s err=%v", uid, token, err)
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusForbidden)
                _ = json.NewEncoder(w).Encode(map[string]interface{}{"code": "token_user_mismatch", "message": http.StatusText(http.StatusForbidden)})
                return
            }
            logx.Errorf("download unknown error: uid=%s token=%s err=%v", uid, token, err)
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusNotFound)
            _ = json.NewEncoder(w).Encode(map[string]interface{}{"code": "unknown", "message": http.StatusText(http.StatusNotFound)})
            return
        }
        // find file path (stored as relative path), and resolve to absolute under storage base
        var rel string
        _ = svcCtx.DB.Conn.QueryRowCtx(r.Context(), &rel, "SELECT file_path FROM models WHERE id=?", mid)
        if rel == "" {
            logx.Errorf("model file path empty: model_id=%d", mid)
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusNotFound)
            _ = json.NewEncoder(w).Encode(map[string]interface{}{"code": "file_missing", "message": http.StatusText(http.StatusNotFound)})
            return
        }
        base := filepath.Clean(svcCtx.Config.StorageBasePath)
        abs := filepath.Clean(filepath.Join(base, filepath.FromSlash(rel)))
        if rp, err := filepath.Rel(base, abs); err != nil || strings.HasPrefix(rp, "..") {
            logx.Errorf("path outside base: abs=%s base=%s", abs, base)
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusForbidden)
            _ = json.NewEncoder(w).Encode(map[string]interface{}{"code": "path_outside_base", "message": http.StatusText(http.StatusForbidden)})
            return
        }
        f, err := os.Open(abs)
        if err != nil {
            logx.Errorf("file open failed: abs=%s err=%v", abs, err)
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusNotFound)
            _ = json.NewEncoder(w).Encode(map[string]interface{}{"code": "file_open_failed", "message": http.StatusText(http.StatusNotFound), "file": rel})
            return
        }
        defer f.Close()
        name := filepath.Base(abs)
        // hint: let browser download with filename
        w.Header().Set("Content-Disposition", "attachment; filename=\""+name+"\"")
        http.ServeContent(w, r, name, time.Now(), f)
    }
}
