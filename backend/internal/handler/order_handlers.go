package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"backend/internal/middleware"
	"backend/internal/repo"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type CreateOrderReq struct {
	Items []int64 `json:"items"`
}

func createOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		sess, ok := middleware.GetSessionFromCtx(r.Context())
		if !ok {
			httpx.ErrorCtx(r.Context(), w, http.ErrNoCookie)
			return
		}
		userId := sess.GetIdentity().Id
		or := repo.NewOrderRepo(svcCtx.DB.Conn)
		oid, err := or.Create(r.Context(), userId, req.Items)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"order_id": oid})
	}
}

func payOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		oid, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		or := repo.NewOrderRepo(svcCtx.DB.Conn)
		o, err := or.Get(r.Context(), oid)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		b := make([]byte, 16)
		_, _ = rand.Read(b)
		idemp := hex.EncodeToString(b)
		pid, err := or.CreatePayment(r.Context(), o.Id, o.TotalCents, idemp)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
			"payment_id":         pid,
			"idempotency_key":    idemp,
			"mock_callback_hint": "/api/payments/callback",
		})
	}
}

type MockCallbackReq struct {
	OrderId        int64  `json:"order_id"`
	Status         string `json:"status"` // succeeded|failed
	IdempotencyKey string `json:"idempotency_key"`
}

func mockCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MockCallbackReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if req.Status != "succeeded" && req.Status != "failed" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		or := repo.NewOrderRepo(svcCtx.DB.Conn)
		if req.Status == "succeeded" {
			cbid := strconv.FormatInt(time.Now().UnixNano(), 10)
			if err := or.ApplyPaymentSucceeded(r.Context(), req.OrderId, cbid, req.IdempotencyKey); err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"ok": true})
	}
}

func listOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, ok := middleware.GetSessionFromCtx(r.Context())
		if !ok {
			httpx.ErrorCtx(r.Context(), w, http.ErrNoCookie)
			return
		}
		userId := sess.GetIdentity().Id
		or := repo.NewOrderRepo(svcCtx.DB.Conn)
		rows, err := or.ListByUser(r.Context(), userId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": rows})
	}
}

func refundOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		oid, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		or := repo.NewOrderRepo(svcCtx.DB.Conn)
		if err := or.ApplyRefund(r.Context(), oid); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"ok": true})
	}
}
