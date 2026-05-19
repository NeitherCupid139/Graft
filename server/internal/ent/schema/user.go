package schema

import (
	"entgo.io/ent"

	userschema "graft/server/plugins/user/ent/schema"
)

// User 保留 internal/ent 的兼容引用面，真正 schema 真值由 user 插件拥有。
type User struct {
	ent.Schema
}

// Fields 转发到 user 插件拥有的 schema 真值。
func (User) Fields() []ent.Field {
	return userschema.User{}.Fields()
}

// Edges 转发到 user 插件拥有的 schema 真值。
func (User) Edges() []ent.Edge {
	return userschema.User{}.Edges()
}
