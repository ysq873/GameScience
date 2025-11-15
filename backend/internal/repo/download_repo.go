package repo

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "time"

    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DownloadRepo struct{ db sqlx.SqlConn }

func NewDownloadRepo(db sqlx.SqlConn) *DownloadRepo { return &DownloadRepo{db: db} }

func (r *DownloadRepo) CreateToken(ctx context.Context, userId string, modelId int64, ttl time.Duration) (string, error) {
    b := make([]byte, 16); _, _ = rand.Read(b)
    token := hex.EncodeToString(b)
    expires := time.Now().Add(ttl)
    _, err := r.db.ExecCtx(ctx, "INSERT INTO download_tokens(token,user_id,model_id,expires_at,one_time) VALUES(?,?,?,?,1)", token, userId, modelId, expires)
    return token, err
}

var ErrTokenExpired = errors.New("download token expired")
var ErrTokenUserMismatch = errors.New("download token user mismatch")

func (r *DownloadRepo) ConsumeToken(ctx context.Context, token string, userId string) (int64, error) {
    var mid int64
    err := r.db.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
        var row struct {
            UserId   string    `db:"user_id"`
            ModelId  int64     `db:"model_id"`
            Expires  time.Time `db:"expires_at"`
        }
        if err := session.QueryRowCtx(ctx, &row, "SELECT user_id, model_id, expires_at FROM download_tokens WHERE token=? FOR UPDATE", token); err != nil { return err }
        if row.UserId != userId { return ErrTokenUserMismatch }
        if time.Now().After(row.Expires) { return ErrTokenExpired }
        mid = row.ModelId
        _, err := session.ExecCtx(ctx, "DELETE FROM download_tokens WHERE token=?", token)
        return err
    })
    return mid, err
}
