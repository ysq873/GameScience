package repo

import (
    "context"
    "database/sql"

    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type FavoriteRepo struct{ db sqlx.SqlConn }

func NewFavoriteRepo(db sqlx.SqlConn) *FavoriteRepo { return &FavoriteRepo{db: db} }

func (r *FavoriteRepo) Add(ctx context.Context, userId string, modelId int64) error {
    _, err := r.db.ExecCtx(ctx, "INSERT IGNORE INTO favorites(user_id,model_id) VALUES(?,?)", userId, modelId)
    return err
}

func (r *FavoriteRepo) Remove(ctx context.Context, userId string, modelId int64) error {
    _, err := r.db.ExecCtx(ctx, "DELETE FROM favorites WHERE user_id=? AND model_id=?", userId, modelId)
    return err
}

type FavModel struct {
    Id    int64
    Title string
    Price int64
    Cover sql.NullString
}

func (r *FavoriteRepo) ListByUser(ctx context.Context, userId string) ([]struct{ ModelId int64 }, error) {
    var rows []struct{ ModelId int64 }
    err := r.db.QueryRowsCtx(ctx, &rows, "SELECT model_id FROM favorites WHERE user_id=? ORDER BY created_at DESC", userId)
    return rows, err
}

func (r *FavoriteRepo) ListModelInfosByUser(ctx context.Context, userId string) ([]FavModel, error) {
    var rows []FavModel
    err := r.db.QueryRowsCtx(ctx, &rows, "SELECT m.id, m.title, m.price_cents, m.cover_url FROM favorites f JOIN models m ON m.id=f.model_id WHERE f.user_id=? ORDER BY f.created_at DESC", userId)
    return rows, err
}