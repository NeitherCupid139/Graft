# entstore

## 用途

`entstore` 提供 `server/internal/store` 契约的 Ent 实现。

## 边界

* Ent client 只在包内出现。
* 对外只返回 `store` 包定义的 DTO 和错误语义。
* 新增查询时优先复用稳定仓储接口，不把 Ent 生成类型直接泄漏给模块。

## 主要入口

* `factory.go`：仓储装配入口
* `user_repository.go`：用户资料读取实现
* `auth_repository.go`：认证口令与 refresh session 实现
* `rbac_repository.go`：角色与权限解析实现

## 维护提示

当 schema 字段或关联发生变化时，需要同步检查 DTO 映射、最小仓储边界以及 Atlas 迁移是否仍保持一致。
