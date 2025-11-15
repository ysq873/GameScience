package handler

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"backend/internal/middleware"
	"backend/internal/svc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func uploadModelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, ok := middleware.GetSessionFromCtx(r.Context())
        if !ok {
            http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
            return
        }
		ownerId := sess.GetIdentity().Id
        if ownerId == "" {
            http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
            return
        }
		_ = r.ParseMultipartForm(64 << 20)
		title := r.FormValue("title")
		description := r.FormValue("description")
		priceStr := r.FormValue("price_cents")
		if title == "" {
			httpx.ErrorCtx(r.Context(), w, errors.New("标题必填"))
			return
		}
		if priceStr == "" {
			httpx.ErrorCtx(r.Context(), w, errors.New("价格(分)必填"))
			return
		}
		price, err := strconv.Atoi(priceStr)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errors.New("价格格式不正确，需整数"))
			return
		}

		modelRel, err := saveFileFromRequest(svcCtx, r, "file", "models")
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if modelRel == "" {
			httpx.ErrorCtx(r.Context(), w, errors.New("模型文件(file)必填"))
			return
		}
		coverRel, _ := saveFileFromRequest(svcCtx, r, "cover", "images")

		res, err := svcCtx.DB.Conn.ExecCtx(r.Context(), "INSERT INTO models(owner_id,title,description,price_cents,cover_url,file_path,status) VALUES(?,?,?,?,?,?,?)",
			ownerId, title, description, price, coverRel, modelRel, 0)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		id, _ := res.LastInsertId()
		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"id": id, "cover_url": coverRel, "file_path": modelRel})
	}
}

func saveFileFromRequest(svcCtx *svc.ServiceContext, r *http.Request, field string, subdir string) (string, error) {
	file, header, err := r.FormFile(field)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return "", nil
		}
		return "", err
	}
	defer file.Close()
	ext := filepath.Ext(header.Filename)
	name := uuid.NewString() + ext
	rel := filepath.Join(subdir, name)
	abs := filepath.Clean(filepath.Join(svcCtx.Config.StorageBasePath, rel))
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return "", err
	}
	dst, err := os.Create(abs)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}
	return rel, nil
}
