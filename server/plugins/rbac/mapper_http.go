package rbac

import (
	"math"
	"time"

	generated "graft/server/internal/contract/openapi/generated"
	rbacstore "graft/server/plugins/rbac/store"
)

func toRoleListResponse(roles []rbacstore.Role) generated.RoleListResponse {
	// Match the generated outer response item shape exactly; the anonymous type keeps the generated `Id` field name.
	items := make([]struct {
		Builtin         bool    `json:"builtin"`
		Description     *string `json:"description,omitempty"`
		Display         string  `json:"display"`
		Id              int64   `json:"id"` //nolint:revive // Must match the generated anonymous response field name.
		Name            string  `json:"name"`
		PermissionCount int     `json:"permission_count"`
		UpdatedAt       string  `json:"updated_at"`
		UserCount       int     `json:"user_count"`
	}, 0, len(roles))
	for _, role := range roles {
		item := toRoleListItem(role)
		items = append(items, struct {
			Builtin         bool    `json:"builtin"`
			Description     *string `json:"description,omitempty"`
			Display         string  `json:"display"`
			Id              int64   `json:"id"` //nolint:revive // Must match the generated anonymous response field name.
			Name            string  `json:"name"`
			PermissionCount int     `json:"permission_count"`
			UpdatedAt       string  `json:"updated_at"`
			UserCount       int     `json:"user_count"`
		}{
			Builtin:         item.Builtin,
			Description:     item.Description,
			Display:         item.Display,
			Id:              item.Id,
			Name:            item.Name,
			PermissionCount: item.PermissionCount,
			UpdatedAt:       item.UpdatedAt,
			UserCount:       item.UserCount,
		})
	}

	return generated.RoleListResponse{Items: items}
}

func toRoleListItem(role rbacstore.Role) generated.RoleListItem {
	return generated.RoleListItem{
		Id:              mustConvertGeneratedID(role.ID, "rbac role id"),
		Name:            role.Name,
		Display:         role.Display,
		Description:     role.Description,
		Builtin:         role.Builtin,
		UpdatedAt:       role.UpdatedAt.UTC().Format(time.RFC3339),
		PermissionCount: role.PermissionCount,
		UserCount:       role.UserCount,
	}
}

func toRolePermissionBindingResponse(bindings []rbacstore.RolePermissionBinding) generated.RolePermissionBindingResponse {
	permissionIDs := make([]int64, 0, len(bindings))
	for _, item := range bindings {
		permissionIDs = append(permissionIDs, mustConvertGeneratedID(item.PermissionID, "rbac permission id"))
	}

	return generated.RolePermissionBindingResponse{PermissionIds: permissionIDs}
}

func toUserRoleBindingResponse(roleIDs []uint64) generated.UserRoleBindingResponse {
	converted := make([]int64, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		converted = append(converted, mustConvertGeneratedID(roleID, "rbac role id"))
	}

	return generated.UserRoleBindingResponse{RoleIds: converted}
}

func toPermissionListResponse(permissions []rbacstore.Permission) generated.PermissionListResponse {
	// Match the generated outer response item shape exactly; the anonymous type keeps the generated `Id` field name.
	items := make([]struct {
		Category         string  `json:"category"`
		Code             string  `json:"code"`
		CreatedAt        string  `json:"created_at"`
		Description      *string `json:"description,omitempty"`
		Display          string  `json:"display"`
		Id               int64   `json:"id"` //nolint:revive // Must match the generated anonymous response field name.
		RoleBindingCount int     `json:"role_binding_count"`
		UpdatedAt        string  `json:"updated_at"`
	}, 0, len(permissions))
	for _, item := range permissions {
		generatedItem := toPermissionListItem(item)
		items = append(items, struct {
			Category         string  `json:"category"`
			Code             string  `json:"code"`
			CreatedAt        string  `json:"created_at"`
			Description      *string `json:"description,omitempty"`
			Display          string  `json:"display"`
			Id               int64   `json:"id"` //nolint:revive // Must match the generated anonymous response field name.
			RoleBindingCount int     `json:"role_binding_count"`
			UpdatedAt        string  `json:"updated_at"`
		}{
			Category:         generatedItem.Category,
			Code:             generatedItem.Code,
			CreatedAt:        generatedItem.CreatedAt,
			Description:      generatedItem.Description,
			Display:          generatedItem.Display,
			Id:               generatedItem.Id,
			RoleBindingCount: generatedItem.RoleBindingCount,
			UpdatedAt:        generatedItem.UpdatedAt,
		})
	}

	return generated.PermissionListResponse{Items: items}
}

func toPermissionListItem(item rbacstore.Permission) generated.PermissionListItem {
	return generated.PermissionListItem{
		Id:               mustConvertGeneratedID(item.ID, "rbac permission id"),
		Code:             item.Code,
		Display:          item.Display,
		Description:      item.Description,
		Category:         item.Category,
		CreatedAt:        item.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:        item.UpdatedAt.UTC().Format(time.RFC3339),
		RoleBindingCount: item.RoleBindingCount,
	}
}

func mustConvertGeneratedID(id uint64, label string) int64 {
	if id > math.MaxInt64 {
		panic(label + " exceeds int64")
	}
	return int64(id)
}
