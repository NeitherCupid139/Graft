package schema

import (
	"time"

	"entgo.io/ent"
	entsql "entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// RolePermission 定义角色与权限点的关联模型。
type RolePermission struct {
	ent.Schema
}

// Annotations 返回 role_permissions 表名映射。
func (RolePermission) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "role_permissions"},
	}
}

// Fields 返回角色权限关联字段定义。
func (RolePermission) Fields() []ent.Field {
	return []ent.Field{
		field.Int("role_id"),
		field.Int("permission_id"),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
	}
}

// Edges 返回角色权限关联的关系定义。
func (RolePermission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", Role.Type).
			Ref("role_permissions").
			Field("role_id").
			Required().
			Unique(),
		edge.From("permission", Permission.Type).
			Ref("role_permissions").
			Field("permission_id").
			Required().
			Unique(),
	}
}

// Indexes 返回角色权限关联的唯一约束与辅助索引。
func (RolePermission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role_id", "permission_id").
			Unique(),
		index.Fields("permission_id"),
	}
}
