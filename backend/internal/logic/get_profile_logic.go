// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"backend/internal/middleware"
	"backend/internal/svc"
	"backend/internal/types"

	ory "github.com/ory/kratos-client-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileLogic) GetProfile() (resp *types.ProfileResp, err error) {
	sess, _ := l.ctx.Value(middleware.CtxKratosSession).(ory.Session)
	id := sess.GetIdentity().Id
	traits := sess.GetIdentity().Traits
	var email string
	if m, ok := traits.(map[string]interface{}); ok {
		if v, ok := m["email"].(string); ok {
			email = v
		}
	}
	r := &types.ProfileResp{Id: id, Email: email}
	return r, nil
}
