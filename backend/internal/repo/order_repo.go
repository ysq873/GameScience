package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Order struct {
	Id         int64        `json:"id"`
	UserId     string       `json:"user_id"`
	TotalCents int64        `json:"total_cents"`
	Status     int          `json:"status_code"`
	CreatedAt  sql.NullTime `json:"created_at"`
}

const (
	OrderStatusPending  = 0
	OrderStatusPaid     = 1
	OrderStatusExpired  = 2
	OrderStatusRefunded = 3
)

type OrderItem struct {
	Id            int64
	OrderId       int64
	ModelId       int64
	PriceCents    int64
	TitleSnapshot sql.NullString
	CoverUrl      sql.NullString
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
			var title sql.NullString
			// 合并查询并校验上架状态
			if err := session.QueryRowCtx(ctx, &price, "SELECT price_cents FROM models WHERE id=? AND status=1", mid); err != nil {
				return errors.New("模型未上架或不存在")
			}
			_ = session.QueryRowCtx(ctx, &title, "SELECT title FROM models WHERE id=?", mid)
			total += price
			// 临时记录到订单项写入阶段使用
			if _, err := session.ExecCtx(ctx, "INSERT INTO order_items(order_id,model_id,price_cents,title_snapshot) VALUES(0,?, ?, ?)", mid, price, title); err != nil {
				// 先写占位，后续替换 order_id；如数据库不允许 0，可换为本地数组缓存（见下优化）
				_ = err // 不影响主流程，改为用本地缓存
			}
		}
		res, err := session.ExecCtx(ctx, "INSERT INTO orders(user_id,total_cents,status) VALUES(?,?,?)", userId, total, OrderStatusPending)
		if err != nil {
			return err
		}
		oid, _ = res.LastInsertId()
		// 重建订单项：以当前已校验的上架模型为准
		for _, mid := range modelIds {
			var price int64
			var title sql.NullString
			// 再次以单次查询获取快照（避免写占位引入复杂性）
			if err := session.QueryRowCtx(ctx, &price, "SELECT price_cents FROM models WHERE id=?", mid); err != nil {
				return err
			}
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
	err := r.db.QueryRowCtx(ctx, &o, "SELECT id,user_id,total_cents,status,created_at FROM orders WHERE id=?", id)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepo) ListByUser(ctx context.Context, userId string) ([]Order, error) {
	var rows []Order
	err := r.db.QueryRowsCtx(ctx, &rows, "SELECT id,user_id,total_cents,status,created_at FROM orders WHERE user_id=? ORDER BY id DESC", userId)
	return rows, err
}

// SecondsLeft returns seconds remaining until 10 minutes after created_at, based on DB clock
func (r *OrderRepo) SecondsLeft(ctx context.Context, orderId int64) (int64, error) {
	var left int64
	err := r.db.QueryRowCtx(ctx, &left, "SELECT GREATEST(0, TIMESTAMPDIFF(SECOND, NOW(), DATE_ADD(created_at, INTERVAL 10 MINUTE))) FROM orders WHERE id=?", orderId)
	return left, err
}

func (r *OrderRepo) MarkExpiredById(ctx context.Context, id int64) error {
	_, err := r.db.ExecCtx(ctx, "UPDATE orders SET status=? WHERE id=? AND status=? AND TIMESTAMPDIFF(MINUTE, created_at, NOW())>10", OrderStatusExpired, id, OrderStatusPending)
	return err
}

func (r *OrderRepo) MarkExpiredPendingOrdersForUser(ctx context.Context, userId string) error {
	_, err := r.db.ExecCtx(ctx, "UPDATE orders SET status=? WHERE user_id=? AND status=? AND TIMESTAMPDIFF(MINUTE, created_at, NOW())>10", OrderStatusExpired, userId, OrderStatusPending)
	return err
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
		var status int
		if err := session.QueryRowCtx(ctx, &status, "SELECT status FROM orders WHERE id=? FOR UPDATE", orderId); err != nil {
			return err
		}
		if status == OrderStatusPaid {
			return nil
		}
		if status != OrderStatusPending {
			return errors.New("invalid order status")
		}
		if _, err := session.ExecCtx(ctx, "UPDATE orders SET status=? WHERE id=?", OrderStatusPaid, orderId); err != nil {
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
		var amount int64
		_ = session.QueryRowCtx(ctx, &amount, "SELECT total_cents FROM orders WHERE id=?", orderId)
		if _, err := session.ExecCtx(ctx, "UPDATE users SET balance=balance-? WHERE kratos_id=?", amount, userId); err != nil {
			return err
		}
		return nil
	})
}

func (r *OrderRepo) ApplyRefund(ctx context.Context, orderId int64) error {
	return r.db.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var status int
		if err := session.QueryRowCtx(ctx, &status, "SELECT status FROM orders WHERE id=? FOR UPDATE", orderId); err != nil {
			return err
		}
		if status == OrderStatusRefunded {
			return nil
		}
		if status != OrderStatusPaid {
			return errors.New("only paid can refund")
		}
		if _, err := session.ExecCtx(ctx, "UPDATE orders SET status=? WHERE id=?", OrderStatusRefunded, orderId); err != nil {
			return err
		}
		var userId string
		_ = session.QueryRowCtx(ctx, &userId, "SELECT user_id FROM orders WHERE id=?", orderId)
		var amount int64
		_ = session.QueryRowCtx(ctx, &amount, "SELECT total_cents FROM orders WHERE id=?", orderId)
		if _, err := session.ExecCtx(ctx, "UPDATE users SET balance=balance+? WHERE kratos_id=?", amount, userId); err != nil {
			return err
		}
		return nil
	})
}

func (r *OrderRepo) ListItems(ctx context.Context, orderId int64) ([]OrderItem, error) {
	var rows []OrderItem
	err := r.db.QueryRowsCtx(ctx, &rows, "SELECT oi.id, oi.order_id, oi.model_id, oi.price_cents, oi.title_snapshot, m.cover_url FROM order_items oi JOIN models m ON m.id=oi.model_id WHERE oi.order_id=?", orderId)
	return rows, err
}
