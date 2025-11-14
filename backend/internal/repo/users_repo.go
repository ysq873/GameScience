package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UsersRepo struct{ db sqlx.SqlConn }

func NewUsersRepo(db sqlx.SqlConn) *UsersRepo { return &UsersRepo{db: db} }

func (r *UsersRepo) GetOrInit(ctx context.Context, kratosId string) (int64, error) {
	var bal int64
	err := r.db.QueryRowCtx(ctx, &bal, "SELECT balance FROM users WHERE kratos_id=?", kratosId)
	if err == sqlx.ErrNotFound {
		_, e := r.db.ExecCtx(ctx, "INSERT INTO users(kratos_id,balance) VALUES(?,0)", kratosId)
		if e != nil {
			return 0, e
		}
		return 0, nil
	}
	return bal, err
}

func (r *UsersRepo) Credit(ctx context.Context, kratosId string, cents int64) error {
	_, err := r.db.ExecCtx(ctx, "UPDATE users SET balance=balance+? WHERE kratos_id=?", cents, kratosId)
	return err
}

func (r *UsersRepo) Debit(ctx context.Context, kratosId string, cents int64) error {
	_, err := r.db.ExecCtx(ctx, "UPDATE users SET balance=balance-? WHERE kratos_id=?", cents, kratosId)
	return err
}

type TxItem struct {
	Id          int64
	OrderId     sql.NullInt64
	Type        string
	AmountCents int64
	Reason      string
	CreatedAt   time.Time
}

func (r *UsersRepo) ListTransactions(ctx context.Context, userId string, page, size int) ([]TxItem, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	offset := (page - 1) * size
	var rows []TxItem
	err := r.db.QueryRowsCtx(ctx, &rows, "SELECT id, order_id, type, amount_cents, reason, created_at FROM account_transactions WHERE user_id=? ORDER BY id DESC LIMIT ? OFFSET ?", userId, size, offset)
	return rows, err
}

func (r *UsersRepo) Recharge(ctx context.Context, userId string, amount int64) (int64, int64, error) {
	var bal int64
	var txId int64
	err := r.db.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		if _, err := s.ExecCtx(ctx, "UPDATE users SET balance_cents=balance_cents+? WHERE user_id=?", amount, userId); err != nil {
			return err
		}
		res, err := s.ExecCtx(ctx, "INSERT INTO account_transactions(user_id,order_id,type,amount_cents,reason) VALUES(?, NULL, 'credit', ?, 'recharge')", userId, amount)
		if err != nil {
			return err
		}
		txId, _ = res.LastInsertId()
		if err := s.QueryRowCtx(ctx, &bal, "SELECT balance_cents FROM users WHERE user_id=?", userId); err != nil {
			return err
		}
		return nil
	})
	return bal, txId, err
}
