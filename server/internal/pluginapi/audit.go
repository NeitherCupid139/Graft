package pluginapi

import "time"

// AuditRecordEventName 是业务插件主动发布审计事件时使用的稳定事件名。
const AuditRecordEventName = "audit.record"

// AuditEvent 描述跨插件可发布的最小审计事件载荷。
//
// 该 DTO 服务于“主动审计”路径：发布方提供明确业务语义，audit 插件负责
// 把事件收敛为稳定持久化记录。
type AuditEvent struct {
	Operator      *CurrentUser
	Action        string
	ResourceType  string
	ResourceID    string
	RequestMethod string
	RequestPath   string
	IP            string
	UserAgent     string
	Success       bool
	ErrorMessage  string
	CreatedAt     time.Time
}
