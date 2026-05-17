import type { RouteRecordRaw } from 'vue-router';

// BootstrapRouteRegistration 描述一个前端模块向壳层声明的 bootstrap 动态路由入口。
//
// 当前阶段只收敛最小运行面：模块显式声明菜单 path、稳定 route name 和页面加载器，
// 壳层再据此把后端 bootstrap 菜单装配成真实动态路由。
export type BootstrapRouteRegistration = {
  menuPath: string;
  routeName: string;
  loadPage: RouteRecordRaw['component'];
};
