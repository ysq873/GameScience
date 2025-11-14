// middleware/kratos_session.go
package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	ory "github.com/ory/kratos-client-go"
	"github.com/zeromicro/go-zero/core/logx"
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

		var sessCookie string
		if raw := r.Header.Get("Cookie"); raw != "" {
			if v := extractCookie(raw, "ory_kratos_session"); v != "" {
				sessCookie = v
			}
		}
		if sessCookie == "" {
			if c, err := r.Cookie("ory_kratos_session"); err == nil {
				sessCookie = c.Value
			}
		}
		if sessCookie != "" {
			if s, resp, err := m.cli.FrontendAPI.ToSession(r.Context()).
				Cookie("ory_kratos_session=" + sessCookie).
				Execute(); err == nil && s.GetActive() {
				r = r.WithContext(context.WithValue(r.Context(), CtxKratosSession, s))
				fmt.Println("whoami:", s)
				next(w, r)
				return
			} else {
				if resp != nil {
					logx.Errorf("kratos whoami failed, status=%d", resp.StatusCode)
				} else {
					logx.Errorf("kratos whoami failed, no response")
				}
			}
		} else {
			logx.Errorf("no ory_kratos_session cookie in request")
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

func extractCookie(raw string, name string) string {
	parts := strings.Split(raw, ";")
	for _, p := range parts {
		kv := strings.SplitN(strings.TrimSpace(p), "=", 2)
		if len(kv) == 2 && kv[0] == name {
			return kv[1]
		}
	}
	return ""
}

func GetSessionFromCtx(ctx context.Context) (ory.Session, bool) {
	v := ctx.Value(CtxKratosSession)
	if v == nil {
		return ory.Session{}, false
	}
	if s, ok := v.(ory.Session); ok {
		return s, true
	}
	if ps, ok := v.(*ory.Session); ok && ps != nil {
		return *ps, true
	}
	return ory.Session{}, false
}
