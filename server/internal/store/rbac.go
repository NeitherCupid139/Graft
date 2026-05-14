package store

import (
	"context"
	"time"
)

// Role 表示 RBAC 角色的稳定持久化 DTO。
type Role struct {
	ID          uint64
	Name        string
	Display     string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Permission 表示 RBAC 权限点的稳定持久化 DTO。
type Permission struct {
	ID          uint64
	Code        string
	Display     string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RBACRepository 暴露未来 RBAC 插件所需的最小持久化查询集合。
//
// 该接口当前只承载角色和权限解析路径，后续如需管理类写操作，应在真实插件需求出现时再收敛扩展。
type RBACRepository interface {
	// ListRolesByUserID 返回指定用户当前绑定的全部角色。
	ListRolesByUserID(ctx context.Context, userID uint64) ([]Role, error)

	// ListPermissionsByUserID 返回指定用户经由角色解析得到的全部权限点。
	ListPermissionsByUserID(ctx context.Context, userID uint64) ([]Permission, error)
}
