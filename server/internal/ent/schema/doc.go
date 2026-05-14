// Package schema 定义 Graft 后端手写的 Ent schema 真值。
//
// 这些 schema 只描述平台当前确认的持久化边界；插件和上层服务必须通过
// `server/internal/store` 暴露的稳定 DTO 与仓储接口访问数据，而不是直接依赖 Ent 生成类型。
package schema
