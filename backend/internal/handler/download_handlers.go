package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"backend/internal/middleware"
	"backend/internal/repo"
	"backend/internal/svc"

	ory "github.com/ory/kratos-client-go"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func listPurchasesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, _ := r.Context().Value(middleware.CtxKratosSession).(ory.Session)
		uid := sess.GetIdentity().Id
		var rows []struct{ ModelId int64 }
		err := svcCtx.DB.Conn.QueryRowsCtx(r.Context(), &rows, "SELECT model_id FROM purchases WHERE user_id=? ORDER BY id DESC", uid)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": rows})
	}
}

func generateDownloadTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, _ := r.Context().Value(middleware.CtxKratosSession).(ory.Session)
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
		sess, _ := r.Context().Value(middleware.CtxKratosSession).(ory.Session)
		uid := sess.GetIdentity().Id
		token := r.URL.Query().Get("token")
		dr := repo.NewDownloadRepo(svcCtx.DB.Conn)
		mid, err := dr.ConsumeToken(r.Context(), token, uid)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		// find file path
		var fp string
		_ = svcCtx.DB.Conn.QueryRowCtx(r.Context(), &fp, "SELECT file_path FROM models WHERE id=?", mid)
		if fp == "" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		abs := filepath.Clean(fp)
		f, err := os.Open(abs)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		defer f.Close()
		http.ServeContent(w, r, filepath.Base(abs), time.Now(), f)
	}
}
