package repo

import (
    "context"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type AccountRepo struct{ db sqlx.SqlConn }

func NewAccountRepo(db sqlx.SqlConn) *AccountRepo { return &AccountRepo{db: db} }

func (r *AccountRepo) GetOrInit(ctx context.Context, userId string) (int64, error) {
    var bal int64
    err := r.db.QueryRowCtx(ctx, &bal, "SELECT balance_cents FROM user_accounts WHERE user_id=?", userId)
    if err == sqlx.ErrNotFound {
        _, e := r.db.ExecCtx(ctx, "INSERT INTO user_accounts(user_id,balance_cents) VALUES(?,0)", userId)
        if e != nil { return 0, e }
        return 0, nil
    }
    return bal, err
}

func (r *AccountRepo) Credit(ctx context.Context, userId string, cents int64) error {
    _, err := r.db.ExecCtx(ctx, "UPDATE user_accounts SET balance_cents=balance_cents+? WHERE user_id=?", cents, userId)
    return err
}

func (r *AccountRepo) Debit(ctx context.Context, userId string, cents int64) error {
    _, err := r.db.ExecCtx(ctx, "UPDATE user_accounts SET balance_cents=balance_cents-? WHERE user_id=?", cents, userId)
    return err
}
