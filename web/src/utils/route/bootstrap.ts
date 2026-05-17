import type { RouteRecordRaw } from 'vue-router';

import type { BootstrapMenu } from '@/api/model/authModel';
import { getBootstrapRouteRegistration } from '@/modules';
import type { AppRouteMeta } from '@/utils/types';

import { localizeRouteTitle } from './title';

const bootstrapLayout: RouteRecordRaw['component'] = () => import('@/layouts/index.vue');

// transformBootstrapMenusToRoutes 把后端 bootstrap 菜单快照映射为当前 web 可消费的最小动态路由。
//
// 当前阶段只接入已在 `web` 内存在页面实现的真实菜单项，避免继续沿用 starter demo 菜单树。
export function transformBootstrapMenusToRoutes(menus: BootstrapMenu[]): RouteRecordRaw[] {
  return menus.flatMap((menu) => {
    const registration = getBootstrapRouteRegistration(menu.path);
    if (!registration) {
      return [];
    }

    const routeMeta: AppRouteMeta = {
      title: localizeRouteTitle(menu.title, menu.title_key),
      titleKey: menu.title_key,
      icon: menu.icon,
      single: true,
    };
    const childMeta: AppRouteMeta = {
      hidden: true,
      title: localizeRouteTitle(menu.title, menu.title_key),
      titleKey: menu.title_key,
    };

    const route = {
      path: menu.path,
      component: bootstrapLayout,
      redirect: `${menu.path}/index`,
      name: registration.routeName,
      meta: routeMeta,
      children: [
        {
          path: 'index',
          name: `${registration.routeName}Index`,
          component: registration.loadPage,
          meta: childMeta,
        },
      ],
    };

    return [route as RouteRecordRaw];
  });
}
