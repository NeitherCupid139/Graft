package httpx

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	actorHeader       = "X-Graft-Actor"
	permissionsHeader = "X-Graft-Permissions"
)

// Session 表示 MVP 阶段请求携带的显式身份信息。
//
// 在真实 auth 与 RBAC 插件落地之前，受保护路由通过这些请求头进行可见
// 的后端权限校验，而不是隐式信任前端路由元数据。
type Session struct {
	Actor       string
	Permissions map[string]struct{}
}

// HasPermission 判断当前会话是否拥有所需权限码。
func (s Session) HasPermission(code string) bool {
	if strings.TrimSpace(code) == "" {
		return true
	}

	_, ok := s.Permissions[code]
	return ok
}

// SessionFromRequest 从请求中解析 MVP 阶段的显式会话头。
func SessionFromRequest(request *http.Request) Session {
	session := Session{
		Actor:       strings.TrimSpace(request.Header.Get(actorHeader)),
		Permissions: make(map[string]struct{}),
	}

	for _, raw := range strings.Split(request.Header.Get(permissionsHeader), ",") {
		permission := strings.TrimSpace(raw)
		if permission == "" {
			continue
		}
		session.Permissions[permission] = struct{}{}
	}

	return session
}

// RequirePermission 执行当前 MVP 阶段的显式权限校验。
func RequirePermission(code string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := SessionFromRequest(ctx.Request)
		if session.Actor == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing request actor",
			})
			return
		}

		if !session.HasPermission(code) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":      "missing permission",
				"permission": code,
			})
			return
		}

		ctx.Next()
	}
}
