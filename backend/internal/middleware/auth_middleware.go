// middleware/kratos_session.go
package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	ory "github.com/ory/kratos-client-go"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type ctxKey string

const CtxKratosSession ctxKey = "kratosSession"

type KratosSessionMiddleware struct {
	cli *ory.APIClient // 指向 Kratos Public API（通常 http://kratos:4433）
}

func NewKratosSessionMiddleware(kratosPublicURL string) *KratosSessionMiddleware {
	cfg := ory.NewConfiguration()
	cfg.Servers = ory.ServerConfigurations{{URL: kratosPublicURL}}
	return &KratosSessionMiddleware{cli: ory.NewAPIClient(cfg)}
}

// 兼容 ServiceContext 的构造函数命名
func NewAuthMiddleware(kratosPublicURL string) *KratosSessionMiddleware {
	return NewKratosSessionMiddleware(kratosPublicURL)
}

func (m *KratosSessionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1) API 流：X-Session-Token / Bearer
		if tok := tokenFromHeaders(r); tok != "" {
			if s, _, err := m.cli.FrontendAPI.ToSession(r.Context()).
				XSessionToken(tok). // 部分 SDK 版本提供该 setter（对应请求头 X-Session-Token）
				Execute(); err == nil && s.GetActive() {
				r = r.WithContext(context.WithValue(r.Context(), CtxKratosSession, s))
				next(w, r)
				return
			}
			// 若你的 SDK 没有 XSessionToken()，也可以自己发起 HTTP 请求设置该头部。
		}

		// 2) 浏览器流：ory_kratos_session Cookie
		if c, err := r.Cookie("ory_kratos_session"); err == nil {
			if s, _, err := m.cli.FrontendAPI.ToSession(r.Context()).
				Cookie("ory_kratos_session=" + c.Value).
				Execute(); err == nil && s.GetActive() {
				r = r.WithContext(context.WithValue(r.Context(), CtxKratosSession, s))
				next(w, r)
				return
			}
		}

		// 无有效会话
		httpx.ErrorCtx(r.Context(), w, errors.New(http.StatusText(http.StatusUnauthorized)))
	}
}

func tokenFromHeaders(r *http.Request) string {
	if v := r.Header.Get("X-Session-Token"); v != "" {
		return v
	}
	if a := r.Header.Get("Authorization"); strings.HasPrefix(strings.ToLower(a), "bearer ") {
		return strings.TrimSpace(a[len("bearer "):])
	}
	return ""
}
