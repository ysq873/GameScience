package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/repo"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type UpdateStatusReq struct {
	Status int `json:"status"`
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
		var statusInt int
		{
			var m map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
			v, ok := m["status"]
			if !ok {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			switch t := v.(type) {
			case float64:
				statusInt = int(t)
			case string:
				switch t {
				case "listed":
					statusInt = 1
				case "delisted":
					statusInt = 2
				case "draft":
					statusInt = 0
				default:
					n, err := strconv.Atoi(t)
					if err != nil {
						http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
						return
					}
					statusInt = n
				}
			default:
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		}
		if statusInt != 0 && statusInt != 1 && statusInt != 2 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		repoM := repo.NewModelRepo(svcCtx.DB.Conn)
		if err := repoM.UpdateStatus(r.Context(), id, strconv.Itoa(statusInt)); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"ok": true})
	}
}
