package middleware

import (
	"context"
	"net/http"
)

type AuthMiddleware struct {
	kratosAdminURL string
}

func NewAuthMiddleware(kratosAdminURL string) *AuthMiddleware {
	return &AuthMiddleware{
		kratosAdminURL: kratosAdminURL,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求中获取 session cookie
		sessionCookie, err := r.Cookie("ory_kratos_session")
		if err != nil {
			http.Error(w, "未认证", http.StatusUnauthorized)
			return
		}

		// 调用 Kratos 管理 API 验证 session
		isValid, identity, err := m.validateSession(r.Context(), sessionCookie.Value)
		if err != nil || !isValid {
			http.Error(w, "会话无效", http.StatusUnauthorized)
			return
		}

		// 将用户身份信息添加到上下文
		ctx := context.WithValue(r.Context(), "userIdentity", identity)
		next(w, r.WithContext(ctx))
	}
}

func (m *AuthMiddleware) validateSession(ctx context.Context, sessionToken string) (bool, map[string]interface{}, error) {
	// 调用 Kratos 的 /sessions/{session} 端点验证 session
	// 这里需要实现 HTTP 客户端调用 Kratos 管理 API
	// 返回 session 是否有效和用户身份信息

	// 简化实现示例：
	client := &http.Client{}
	req, _ := http.NewRequest("GET", m.kratosAdminURL+"/admin/sessions/"+sessionToken, nil)
	resp, err := client.Do(req)
	if err != nil {
		return false, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil, nil
	}

	// 解析响应获取身份信息
	// var sessionResp SessionResponse
	// json.NewDecoder(resp.Body).Decode(&sessionResp)

	return true, map[string]interface{}{"id": "user-id"}, nil
}
