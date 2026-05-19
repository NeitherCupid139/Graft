// Package schema 定义 Graft 后端当前仍由 internal/ent 持有的手写 Ent schema，
// 以及为既有生成包保留的最小兼容引用面。
//
// plugin-owned schema 真值应收敛到各自插件目录；这里只保留 core-owned 定义或
// 过渡期兼容别名，避免把共享 generated package 重新当作业务真值来源。
package schema
