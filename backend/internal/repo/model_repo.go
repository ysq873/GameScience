package repo

import (
    "context"
    "database/sql"

    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Model struct {
    Id         int64
    OwnerId    string
    Title      string
    Description string
    PriceCents int64
    CoverUrl   sql.NullString
    FilePath   string
    Status     string
}

type ModelRepo struct{
    db sqlx.SqlConn
}

func NewModelRepo(db sqlx.SqlConn) *ModelRepo { return &ModelRepo{db: db} }

func (r *ModelRepo) ListListed(ctx context.Context, page, size int) ([]Model, error) {
    if size <= 0 { size = 20 }
    if page <= 0 { page = 1 }
    offset := (page-1)*size
    query := "SELECT id, owner_id, title, description, price_cents, cover_url, file_path, status FROM models WHERE status='listed' ORDER BY id DESC LIMIT ? OFFSET ?"
    var rows []Model
    err := r.db.QueryRowsCtx(ctx, &rows, query, size, offset)
    return rows, err
}

func (r *ModelRepo) GetByID(ctx context.Context, id int64) (*Model, error) {
    var m Model
    err := r.db.QueryRowCtx(ctx, &m, "SELECT id, owner_id, title, description, price_cents, cover_url, file_path, status FROM models WHERE id=?", id)
    if err != nil { return nil, err }
    return &m, nil
}

func (r *ModelRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
    _, err := r.db.ExecCtx(ctx, "UPDATE models SET status=? WHERE id=?", status, id)
    return err
}
