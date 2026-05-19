package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"

	userschema "graft/server/plugins/user/ent/schema"
)

// RefreshSession 保留 internal/ent 的兼容引用面，真正 schema 真值由 user 插件拥有。
type RefreshSession struct {
	ent.Schema
}

// Annotations 转发到 user 插件拥有的 schema 真值。
func (RefreshSession) Annotations() []schema.Annotation {
	return userschema.RefreshSession{}.Annotations()
}

// Fields 转发到 user 插件拥有的 schema 真值。
func (RefreshSession) Fields() []ent.Field {
	return userschema.RefreshSession{}.Fields()
}

// Edges 转发到 user 插件拥有的 schema 真值。
func (RefreshSession) Edges() []ent.Edge {
	return userschema.RefreshSession{}.Edges()
}

// Indexes 转发到 user 插件拥有的 schema 真值。
func (RefreshSession) Indexes() []ent.Index {
	return userschema.RefreshSession{}.Indexes()
}
