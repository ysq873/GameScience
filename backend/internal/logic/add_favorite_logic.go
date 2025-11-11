// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFavoriteLogic {
	return &AddFavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddFavoriteLogic) AddFavorite(req *types.AddFavoriteReq) (resp *types.ProfileResp, err error) {
	// todo: add your logic here and delete this line

	return
}
