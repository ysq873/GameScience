package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Order struct {
	Id         int64
	UserId     string
	TotalCents int64
	Status     string
}

type OrderItem struct {
	Id            int64
	OrderId       int64
	ModelId       int64
	PriceCents    int64
	TitleSnapshot sql.NullString
}

type Payment struct {
	Id             int64
	OrderId        int64
	AmountCents    int64
	Status         string
	Gateway        string
	IdempotencyKey sql.NullString
	CallbackId     sql.NullString
}

type OrderRepo struct{ db sqlx.SqlConn }

func NewOrderRepo(db sqlx.SqlConn) *OrderRepo { return &OrderRepo{db: db} }

func (r *OrderRepo) Create(ctx context.Context, userId string, modelIds []int64) (int64, error) {
	var oid int64
	err := r.db.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var total int64
		for _, mid := range modelIds {
			var price int64
			if err := session.QueryRowCtx(ctx, &price, "SELECT price_cents FROM models WHERE id=? AND status='listed'", mid); err != nil {
				return err
			}
			total += price
		}
		res, err := session.ExecCtx(ctx, "INSERT INTO orders(user_id,total_cents,status) VALUES(?,?,'pending')", userId, total)
		if err != nil {
			return err
		}
		oid, _ = res.LastInsertId()
		for _, mid := range modelIds {
			var price int64
			var title sql.NullString
			_ = session.QueryRowCtx(ctx, &price, "SELECT price_cents FROM models WHERE id=?", mid)
			_ = session.QueryRowCtx(ctx, &title, "SELECT title FROM models WHERE id=?", mid)
			if _, err := session.ExecCtx(ctx, "INSERT INTO order_items(order_id,model_id,price_cents,title_snapshot) VALUES(?,?,?,?)", oid, mid, price, title); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return oid, nil
}

func (r *OrderRepo) Get(ctx context.Context, id int64) (*Order, error) {
	var o Order
	err := r.db.QueryRowCtx(ctx, &o, "SELECT id,user_id,total_cents,status FROM orders WHERE id=?", id)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepo) ListByUser(ctx context.Context, userId string) ([]Order, error) {
	var rows []Order
	err := r.db.QueryRowsCtx(ctx, &rows, "SELECT id,user_id,total_cents,status FROM orders WHERE user_id=? ORDER BY id DESC", userId)
	return rows, err
}

func (r *OrderRepo) CreatePayment(ctx context.Context, orderId int64, amount int64, idempKey string) (int64, error) {
	res, err := r.db.ExecCtx(ctx, "INSERT INTO payments(order_id,amount_cents,status,gateway,idempotency_key) VALUES(?,?, 'initiated','mock', ?)", orderId, amount, idempKey)
	if err != nil {
		return 0, err
	}
	pid, _ := res.LastInsertId()
	return pid, nil
}

func (r *OrderRepo) ApplyPaymentSucceeded(ctx context.Context, orderId int64, callbackId string, idempKey string) error {
	return r.db.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var status string
		if err := session.QueryRowCtx(ctx, &status, "SELECT status FROM orders WHERE id=? FOR UPDATE", orderId); err != nil {
			return err
		}
		if status == "paid" {
			return nil
		}
		if status != "pending" {
			return errors.New("invalid order status")
		}
		if _, err := session.ExecCtx(ctx, "UPDATE orders SET status='paid' WHERE id=?", orderId); err != nil {
			return err
		}
		if _, err := session.ExecCtx(ctx, "UPDATE payments SET status='succeeded', callback_id=? WHERE order_id=? AND idempotency_key=?", callbackId, orderId, idempKey); err != nil {
			return err
		}
		var userId string
		_ = session.QueryRowCtx(ctx, &userId, "SELECT user_id FROM orders WHERE id=?", orderId)
		var items []struct{ ModelId int64 }
		if err := session.QueryRowsCtx(ctx, &items, "SELECT model_id FROM order_items WHERE order_id=?", orderId); err != nil {
			return err
		}
		for _, it := range items {
			_, _ = session.ExecCtx(ctx, "INSERT IGNORE INTO purchases(user_id,model_id,order_id) VALUES(?,?,?)", userId, it.ModelId, orderId)
		}
		return nil
	})
}

func (r *OrderRepo) ApplyRefund(ctx context.Context, orderId int64) error {
	return r.db.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var status string
		if err := session.QueryRowCtx(ctx, &status, "SELECT status FROM orders WHERE id=? FOR UPDATE", orderId); err != nil {
			return err
		}
		if status == "refunded" {
			return nil
		}
		if status != "paid" {
			return errors.New("only paid can refund")
		}
		if _, err := session.ExecCtx(ctx, "UPDATE orders SET status='refunded' WHERE id=?", orderId); err != nil {
			return err
		}
		return nil
	})
}
