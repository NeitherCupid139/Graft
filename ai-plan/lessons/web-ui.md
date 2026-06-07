# Web UI Lessons

## LESSON-WEB-UI-PROTECTED-STATE-001：系统保护状态不应伪装成错误告警

- Status: active
- Level: L1
- Applies to:
  - `web` management pages with builtin/system/protected records
  - readonly drawers, dialogs, and action menus
  - `list-form-detail` page type
- Source:
  - role management system-role UX remediation
  - user feedback that builtin roles are normal business state, not warning/error conditions
- Problem:
  系统内置角色、只读权限、受保护配置这类状态是平台的正常保护模型。若用橙色 warning 块、置灰可写按钮或隐藏原因来表达，
  用户会误以为数据异常或权限系统出错，也不知道下一步能做什么。
- Correct pattern:
  对受保护但正常的业务状态，应使用 info、neutral 或 primary-light 语义，文案明确说明“这是正常限制，不是异常”。
  操作模型应从可写操作切换为只读操作，例如“查看权限”替代禁用的“分配权限”；更多菜单只暴露可执行动作，例如详情、
  查看、复制为自定义对象。只读弹窗应允许搜索、展开和查看说明，但隐藏保存按钮并保留关闭动作。
- Anti-pattern:
  - 用 warning/error 样式表达正常系统保护状态
  - 保留可写按钮但简单 disabled，且没有说明原因
  - 对系统内置对象暴露编辑、删除、分配权限等不可执行操作
  - 只在前端隐藏危险操作，却不确认后端 authority 有拒绝路径
- Enforcement:
  修改带 builtin/system/protected 字段的管理页时，检查表格操作、详情抽屉、弹窗标题、提示块和 footer 是否按只读/保护态
  切换；确认文案在模块 i18n 中；确认后端或 contract authority 已拒绝不可执行写操作。
- Promotion:
  - AGENTS.md: no
  - Design doc: no
- Related:
  - `web/src/modules/rbac/pages/index.vue`
  - `web/src/shared/components/assignment/AssignmentFooter.vue`
- Updated at:
  2026-06-07

## LESSON-WEB-UI-EMPTY-STATE-001：表格空状态不应做成小灰色卡片

- Status: active
- Level: L3
- Applies to:
  - `web` table/list management pages
  - `list-form-detail` page type
  - TDesign Vue Next table empty states
- Source:
  - 用户管理 / 角色管理空状态修复
  - user feedback on wrong empty-state implementation direction
- Problem:
  AI 曾在用户管理、角色管理这类表格页中，把“暂无数据”实现成表格中间的小灰色卡片，视觉上像浮层或占位块。
  这种做法不符合 TDesign 管理页空态模式，也容易与表格容器、分页和暗色主题产生割裂。
- Correct pattern:
  对于 table/list management 页面，应优先使用 `t-empty` 组件或 `t-table` 的 `empty` 插槽。空状态应位于 table body
  中央，保留表头和分页结构，颜色与边框必须使用 TDesign token。创建型管理页可提供主操作按钮；当筛选条件激活时，可同时
  提供“清空筛选”操作。
- Anti-pattern:
  - 手写 `empty-card` 或 `empty-box`
  - 在表格中间放一个小灰色盒子
  - 硬编码 `#f5f5f5`、`#fff`、`#000`
  - 让空状态挤乱分页
  - 让空状态与表格 header/body/footer 结构割裂
  - 让暗色主题下的空状态背景或文案失去可读性
- Enforcement:
  实现或修改 table/list management 页面时，必须检查 empty state 是否使用 `t-empty` 或 table empty slot，分页是否稳定，
  表头/空态/分页结构是否连续，以及颜色是否全部来自 TDesign token。
- Promotion:
  - AGENTS.md: yes
  - Design doc: yes
- Related:
  - `web/AGENTS.md`
  - `ai-plan/design/graft-design-system/list-form-detail.md`
- Updated at:
  2026-05-22

## LESSON-WEB-UI-LOG-AUDIT-001：高级查询列表页必须优先抽通用查询结构

- Status: active
- Level: L2
- Applies to:
  - `query-builder-list-detail` page type
  - `web` log and audit pages
  - `log-audit` page type
  - access-log, app-log, audit logs, and future field-heavy query pages
- Source:
  - app-log page remediation after user feedback that the page did not follow the access-log page pattern
  - duplicate local implementations discovered while aligning app-log with access-log
  - query-list refactor after user feedback that the abstraction should serve future field-heavy query pages, not only logs
- Problem:
  字段多、筛选复杂、需要分页表格和详情抽屉的页面若按单页临时实现，很容易出现筛选器、表格排序、详情抽屉、深链参数、
  空态、错误提示和交互文案不一致。这种分叉会让相同类型页面看起来像不同产品，也会在严格重复代码检查下产生维护成本。
- Correct pattern:
  新增或重做字段密集查询页时，先声明页面类型为 `query-builder-list-detail`；日志审计只是 `log-audit` 变体。页面壳、
  筛选构建器、分页表格、列设置抽屉、列表错误提示和通用交互应优先沉淀到 `web/src/shared/components/query-list`，
  由页面只提供领域字段、API 查询、深链语义、详情组件和展示文案。
- Anti-pattern:
  - 在每个字段密集查询页手写一套筛选器、表格、列设置、详情抽屉和错误反馈
  - 只复制访问日志页面的视觉结果，却保留不同的数据流、外壳结构和交互语义
  - 用本地兼容映射掩盖后端契约缺少排序、筛选或分页字段的问题
  - 为通过重复代码检查而做无语义的改名或拆行
- Enforcement:
  修改字段密集查询页时，检查页面类型是否为 `query-builder-list-detail` 或其变体，并确认页面壳、筛选器、分页表格、
  列设置和错误态是否复用 `shared/components/query-list`。业务字段、API 查询、URL deep-link 和详情内容仍留在模块内。
  用 `bun run dupcode:check`、相关 Vitest 用例和 `bun run check` 验证没有重新分叉。
- Promotion:
  - AGENTS.md: no
  - Design doc: no
- Related:
  - `web/src/shared/components/query-list`
  - `web/src/modules/access-log/pages/list/index.vue`
  - `web/src/modules/app-log/pages/list/index.vue`
  - `web/src/modules/audit/pages/logs/index.vue`
- Updated at:
  2026-06-04

## LESSON-WEB-UI-PAGE-CONTAINER-001：后台页面容器应统一复用共享容器与宽度变量策略

- Status: active
- Level: L2
- Applies to:
  - `web` management pages
  - shared page containers such as `ManagementPageContent`
  - large-screen layout tuning and visual-centering fixes
- Source:
  - cross-page layout consistency review for access control, user, role, and monitor pages
  - repeated left/right whitespace and centering drift discussions
- Problem:
  后台页面在大屏下容易出现左右留白、宽度策略不一致和视觉重心漂移。若每个页面单独修宽度或偏移，长期会导致容器策略失控。
- Correct pattern:
  后台页面应优先复用共享页面容器与宽度变量策略，例如统一的 `max-width`、`width-ratio`、`min-padding` 和
  `margin-inline: auto`。管理页可以通过变量覆盖获得更宽的内容面，但不能破坏整体居中。排查时应优先查看 DOM/CSS
  计算结果，而不是只凭截图主观修偏移。
- Anti-pattern:
  - 单页写死 `margin-left`
  - 用 `transform` 修视觉偏移
  - 每个页面自行定义长期 `max-width` 策略
  - 忽略滚动条、悬浮工具按钮或容器嵌套对视觉重心的影响
  - 让共享容器和页级容器同时争夺宽度真值
- Enforcement:
  新增或调整后台页面宽度时，先检查是否已有共享容器可复用；若需要覆盖宽度，必须通过变量而不是局部偏移 hack；若出现视觉不居中，
  需要检查实际容器计算宽度与布局约束。
- Promotion:
  - AGENTS.md: no
  - Design doc: yes
- Related:
  - `ai-plan/design/前端视觉设计规范.md`
  - `web/src/shared/components/management/ManagementPageContent.vue`
- Updated at:
  2026-05-22

## LESSON-WEB-UI-DENSITY-TOKEN-001：信息密度切换必须治理 token 消费面

- Status: active
- Level: L2
- Applies to:
  - `web` theme workbench and density presets
  - `web/src/style/**` global layout tokens
  - `web/src/layouts/**`, `web/src/shared/**`, and module page styles
- Source:
  - information-density switch remediation after user feedback that only a small amount of text spacing changed
  - full-web density governance gate added for `web/src/**`
- Problem:
  信息密度 preset 如果只更新 TDesign component size token 或单个 `--graft-theme-density-scale`，而页面继续写死
  `gap`、`padding`、`margin`、`t-space size` 和图表 tooltip 内联间距，设置面板会显示可切换，但真实页面节奏几乎不变。
  这种分叉会让主题工作台变成“看得见的配置、感受不到的体验”。
- Correct pattern:
  信息密度能力必须同时覆盖 source token 和消费面。密度相关布局应使用 `--td-comp-*`、`--graft-density-*` 或
  `calc(...var(--graft-theme-density-scale)...)`，并让共享组件、业务页面、TDesign `Space` 尺寸和图表 tooltip
  内联模板共同响应同一套 density authority。
- Anti-pattern:
  - 只在 store 里生成密度 token，不替换页面里的固定间距
  - 在业务页面继续新增裸 `gap: 16px`、`padding: 12px 14px` 或 `<t-space size="8px">`
  - 用局部 class 名如 `compact` 伪装为全局信息密度响应
  - 把图表 tooltip HTML 字符串排除在密度治理之外
- Enforcement:
  修改前端布局或新增页面时，运行 `bun run density:check` 或完整 `bun run check`。扫描发现的固定密度间距必须改为
  TDesign/Graft density token；只有图标盒、断点、安全区、滚动条、媒体尺寸等非信息密度几何值才允许进入脚本白名单，
  且白名单必须写明具体原因。
- Promotion:
  - AGENTS.md: no
  - Design doc: no
- Related:
  - `web/scripts/check-density-governance.ts`
  - `web/src/store/modules/setting.ts`
  - `web/src/style/layout.less`
- Updated at:
  2026-06-05
