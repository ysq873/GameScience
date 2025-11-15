package handler

import (
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "net/http"
    "strconv"
    "time"
    "strings"

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
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		userId := sess.GetIdentity().Id
        or := repo.NewOrderRepo(svcCtx.DB.Conn)
        if len(req.Items) > 0 {
            qs := make([]string, 0, len(req.Items))
            args := make([]interface{}, 0, len(req.Items)+1)
            args = append(args, userId)
            for _, id := range req.Items { qs = append(qs, "?"); args = append(args, id) }
            var dups []int64
            _ = svcCtx.DB.Conn.QueryRowsCtx(r.Context(), &dups, "SELECT model_id FROM purchases WHERE user_id=? AND model_id IN ("+strings.Join(qs, ",")+")", args...)
            if len(dups) > 0 {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                _ = json.NewEncoder(w).Encode(map[string]interface{}{"code": "already_purchased", "duplicates": dups, "message": "已购买过"})
                return
            }
        }
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
		// 10分钟有效期校验（以DB时间为准，避免时区误差）
		left, _ := repo.NewOrderRepo(svcCtx.DB.Conn).SecondsLeft(r.Context(), o.Id)
		if left <= 0 {
			_ = repo.NewOrderRepo(svcCtx.DB.Conn).MarkExpiredById(r.Context(), o.Id)
			http.Error(w, "订单已过期，请重新下单", http.StatusBadRequest)
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
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		userId := sess.GetIdentity().Id
		or := repo.NewOrderRepo(svcCtx.DB.Conn)
		_ = or.MarkExpiredPendingOrdersForUser(r.Context(), userId)
		rows, err := or.ListByUser(r.Context(), userId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		out := make([]map[string]interface{}, 0, len(rows))
		for _, o := range rows {
			st := "pending"
			if o.Status == repo.OrderStatusPaid {
				st = "paid"
			} else if o.Status == repo.OrderStatusExpired {
				st = "expired"
			} else if o.Status == repo.OrderStatusRefunded {
				st = "refunded"
			}
			m := map[string]interface{}{"id": o.Id, "user_id": o.UserId, "total_cents": o.TotalCents, "status": st, "status_code": o.Status}
			if o.CreatedAt.Valid {
				m["created_at"] = o.CreatedAt.Time.Format(time.RFC3339)
			}
			out = append(out, m)
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": out})
	}
}

func orderDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
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
		sess, ok := middleware.GetSessionFromCtx(r.Context())
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		or := repo.NewOrderRepo(svcCtx.DB.Conn)
		o, err := or.Get(r.Context(), oid)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if o.UserId != sess.GetIdentity().Id {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		items, err := or.ListItems(r.Context(), oid)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		// build response
		resp := map[string]interface{}{
			"id":          o.Id,
			"user_id":     o.UserId,
			"total_cents": o.TotalCents,
			"status_code": o.Status,
		}
		st := "pending"
		if o.Status == repo.OrderStatusPaid {
			st = "paid"
		} else if o.Status == repo.OrderStatusExpired {
			st = "expired"
		} else if o.Status == repo.OrderStatusRefunded {
			st = "refunded"
		}
		resp["status"] = st
		if o.CreatedAt.Valid {
			created := o.CreatedAt.Time
			resp["created_at"] = created.UTC().Format(time.RFC3339)
			expires := created.Add(10 * time.Minute)
			resp["expires_at"] = expires.UTC().Format(time.RFC3339)
			// 以DB计算的秒数为准，避免跨时区误差
			left, _ := repo.NewOrderRepo(svcCtx.DB.Conn).SecondsLeft(r.Context(), o.Id)
			if left < 0 {
				left = 0
			}
			resp["expires_seconds_left"] = int(left)
			if left == 0 && o.Status == repo.OrderStatusPending {
				_ = repo.NewOrderRepo(svcCtx.DB.Conn).MarkExpiredById(r.Context(), o.Id)
				resp["status"] = "expired"
			}
		}
		resp["server_now"] = time.Now().UTC().Format(time.RFC3339)
		its := make([]map[string]interface{}, 0, len(items))
		for _, it := range items {
			title := ""
			if it.TitleSnapshot.Valid {
				title = it.TitleSnapshot.String
			}
			cover := ""
			if it.CoverUrl.Valid {
				cover = it.CoverUrl.String
			}
			its = append(its, map[string]interface{}{"model_id": it.ModelId, "title": title, "price_cents": it.PriceCents, "cover_url": cover, "quantity": 1, "subtotal_cents": it.PriceCents})
		}
		resp["items"] = its
		httpx.OkJsonCtx(r.Context(), w, resp)
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
