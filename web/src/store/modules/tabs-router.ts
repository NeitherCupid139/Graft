import { defineStore } from 'pinia';

import { LOCALE } from '@/contracts/i18n/locales';
import { AUTH_ROUTE_NAME } from '@/modules/auth/contract/routes';
import type { TRouterInfo, TTabRouterType } from '@/utils/types';

const homeRoute: Array<TRouterInfo> = [
  {
    path: '/',
    routeIdx: 0,
    title: { [LOCALE.ZH_CN]: '首页', [LOCALE.EN_US]: 'Home' },
    name: 'RootEntry',
    isHome: true,
  },
];

// 不需要做多标签tabs页缓存的列表 值为每个页面对应的name 如 DashboardDetail
// const ignoreCacheRoutes = ['DashboardDetail'];
const ignoreCacheRoutes: string[] = [AUTH_ROUTE_NAME.LOGIN];

function createInitialState(): TTabRouterType {
  return {
    tabRouterList: homeRoute.map((route) => ({ ...route })),
    isRefreshing: false,
  };
}

function shouldKeepTabAlive(route: TRouterInfo) {
  return !route.isHome && !ignoreCacheRoutes.includes(route.name as string) && route.meta?.keepAlive !== false;
}

export const useTabsRouterStore = defineStore('tabsRouter', {
  state: createInitialState,
  getters: {
    tabRouters: (state: TTabRouterType) => state.tabRouterList,
    refreshing: (state: TTabRouterType) => state.isRefreshing,
  },
  actions: {
    startTabRefresh(routeIdx: number) {
      const route = this.tabRouters[routeIdx];
      if (!route) {
        this.isRefreshing = false;
        return;
      }

      this.isRefreshing = true;
      route.isAlive = false;
    },
    finishTabRefresh(routeIdx: number) {
      const route = this.tabRouters[routeIdx];
      if (route) {
        route.isAlive = shouldKeepTabAlive(route);
      }

      this.isRefreshing = false;
    },
    healPersistedState() {
      this.isRefreshing = false;
      this.tabRouterList = this.tabRouters.map((route) => ({
        ...route,
        isAlive: route.isHome ? true : shouldKeepTabAlive(route),
      }));
    },
    // 处理新增
    appendTabRouterList(newRoute: TRouterInfo) {
      // 不要将判断条件newRoute.meta.keepAlive !== false修改为newRoute.meta.keepAlive，starter默认开启保活，所以meta.keepAlive未定义时也需要进行保活，只有显式说明false才禁用保活。
      const needAlive = shouldKeepTabAlive(newRoute);
      if (!this.tabRouters.find((route: TRouterInfo) => route.path === newRoute.path)) {
        this.tabRouterList = this.tabRouterList.concat({ ...newRoute, isAlive: needAlive });
      }
    },
    // 处理关闭当前
    subtractCurrentTabRouter(newRoute: TRouterInfo) {
      const { routeIdx } = newRoute;
      if (routeIdx === undefined) return;
      this.tabRouterList = this.tabRouterList.slice(0, routeIdx).concat(this.tabRouterList.slice(routeIdx + 1));
    },
    // 处理关闭右侧
    subtractTabRouterBehind(newRoute: TRouterInfo) {
      const { routeIdx } = newRoute;
      if (routeIdx === undefined) return;
      const homeIdx: number = this.tabRouters.findIndex((route: TRouterInfo) => route.isHome);
      let tabRouterList: Array<TRouterInfo> = this.tabRouterList.slice(0, routeIdx + 1);
      if (routeIdx < homeIdx) {
        tabRouterList = tabRouterList.concat(homeRoute);
      }
      this.tabRouterList = tabRouterList;
    },
    // 处理关闭左侧
    subtractTabRouterAhead(newRoute: TRouterInfo) {
      const { routeIdx } = newRoute;
      if (routeIdx === undefined) return;
      const homeIdx: number = this.tabRouters.findIndex((route: TRouterInfo) => route.isHome);
      let tabRouterList: Array<TRouterInfo> = this.tabRouterList.slice(routeIdx);
      if (routeIdx > homeIdx) {
        tabRouterList = homeRoute.concat(tabRouterList);
      }
      this.tabRouterList = tabRouterList;
    },
    // 处理关闭其他
    subtractTabRouterOther(newRoute: TRouterInfo) {
      const { routeIdx } = newRoute;
      if (routeIdx === undefined) return;
      const homeIdx: number = this.tabRouters.findIndex((route: TRouterInfo) => route.isHome);
      this.tabRouterList = routeIdx === homeIdx ? homeRoute : homeRoute.concat([this.tabRouterList[routeIdx]]);
    },
    removeTabRouterList() {
      this.tabRouterList = [];
    },
    initTabRouterList(newRoutes: TRouterInfo[]) {
      newRoutes?.forEach((route: TRouterInfo) => this.appendTabRouterList(route));
    },
  },
  persist: {
    pick: ['tabRouterList'],
  },
});
