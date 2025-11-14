package handler

import (
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"

    "backend/internal/svc"
)

func staticHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        rel := r.URL.Query().Get("file")
        if rel == "" {
            http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
            return
        }
        rel = strings.ReplaceAll(rel, "\\", "/")
        base := filepath.Clean(svcCtx.Config.StorageBasePath)
        abs := filepath.Clean(filepath.Join(base, rel))
        if !strings.HasPrefix(abs, base) {
            http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
            return
        }
        f, err := os.Open(abs)
        if err != nil {
            http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
            return
        }
        defer f.Close()
        http.ServeContent(w, r, filepath.Base(abs), time.Now(), f)
    }
}

