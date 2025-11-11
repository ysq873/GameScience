// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"backend/internal/config"
	"backend/internal/middleware"

	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware("TODO: 填写 Kratos 管理 API 地址").Handle,
	}
}
