# 后端主导的 MVP 闭环收敛计划

## 当前阶段概述

- `mvp-extension-path` 已完成最小 `server` / `web` 壳层与恢复材料拆分。
- 下一阶段不再继续横向扩大会话治理或页面数量，而是先由 `server` 补齐 MVP 必要闭环，再让 `web` 收敛到真实契约。

## 下一阶段目标

- 完成 `server` 侧 event bus、audit、scheduler 的最小可执行闭环。
- 稳定当前阶段必须共享的 `auth`、`menu`、`permission`、`locale` 与插件公开 DTO 契约。
- 让 `web` 在 starter 壳层基础上接入真实后端契约，并尽快退出 mock/demo 依赖。

## 实施顺序

1. 稳定 event bus 与相关插件边界，避免后续 audit / scheduler 再次回改 core。
2. 落地最小 audit 路径，明确事件写入、查询边界和权限收敛面。
3. 落地 scheduler plugin 与 cron/runtime 注册闭环，补齐 MVP 级后台任务能力。
4. 冻结当前阶段需要给 `web` 消费的真实契约，包括 `auth`、`menu`、`permission`、`locale` 与必要公开 DTO。
5. 让 `web` 用 starter 壳层挂接真实登录、菜单、权限与本地化契约，清理 mock/demo 依赖。

## MVP 必须完成

- `server` 具备最小 event bus 能力，并可支撑 audit / scheduler 的真实执行路径。
- `audit` 具备最小记录闭环，而不是停留在预留接口。
- `scheduler` 具备最小可注册、可启动、可关闭的任务闭环。
- 共享后端契约足够稳定，`web` 能基于真实接口完成登录、菜单、权限和 locale 收敛。
- `web` 当前阶段完成 starter 壳层收敛，不以新增页面为完成标准。

## 后续递延

- 现有会话治理基线之外的更多 revoke/filter/list 宽度。
- 新页面、新模块和非闭环必需的前端视觉扩张。
- 深度 bundle 优化、主题工作台深化和非关键依赖升级。
- 超出 MVP 闭环所需的更重后台治理或平台化抽象。

## 当前 Web 阶段策略

- `web` 当前只做 starter 壳层收敛 + 真实后端契约挂接。
- 优先接入真实 `auth`、`current user`、动态菜单、权限门禁和本地化错误/文案契约。
- 不以页面扩张为目标；新增页面只能在真实契约已经稳定后再评估。
- 任何临时 mock/demo 入口都应被快速隔离或移除，避免重新形成前端自洽的假闭环。
