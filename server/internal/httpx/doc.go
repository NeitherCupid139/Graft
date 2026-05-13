// Package httpx 管理显式的 HTTP 服务外壳。
//
// 这个包负责 Gin 服务包装、MVP 阶段的后端权限守卫，以及服务启动与
// 关闭语义，避免把这些生命周期细节散落到 core 和业务插件中。
package httpx
