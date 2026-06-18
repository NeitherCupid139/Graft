// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import { describe, expect, it } from 'vitest';
import type { RouteRecordRaw } from 'vue-router';

import { buildDashboardQuickActionLinks } from './sidebar-quick-actions';

function asRouteRecordRaw<T extends object>(route: T) {
  return route as unknown as RouteRecordRaw;
}

describe('buildDashboardQuickActionLinks', () => {
  it('keeps only visible sidebar leaf routes and preserves route title/icon truth', () => {
    const routes = [
      {
        path: '/access-control',
        name: 'BootstrapGroupAccessControl',
        meta: {
          titleKey: 'menu.accessControl',
        },
        children: [
          asRouteRecordRaw({
            path: 'roles',
            name: 'RoleListIndex',
            meta: {
              icon: 'secured',
              orderNo: 20,
              tabTitle: {
                'zh-CN': '访问控制 - 角色管理',
                'en-US': 'Access Control - Roles',
              },
              titleKey: 'rbac.role.list.title',
            },
          }),
          asRouteRecordRaw({
            path: 'permissions',
            name: 'PermissionListIndex',
            meta: {
              hiddenMenu: true,
              orderNo: 30,
              tabTitle: {
                'zh-CN': '访问控制 - 权限管理',
                'en-US': 'Access Control - Permissions',
              },
              titleKey: 'rbac.permission.list.title',
            },
          }),
        ],
      },
      {
        path: '/ops/containers',
        name: 'ContainerList',
        meta: {
          icon: 'layers',
          orderNo: 10,
          single: true,
          tabTitle: {
            'zh-CN': '容器管理',
            'en-US': 'Container Management',
          },
          titleKey: 'container.route.list.title',
        },
        children: [
          asRouteRecordRaw({
            path: 'index',
            name: 'ContainerListIndex',
            meta: {
              hidden: true,
              titleKey: 'container.route.list.title',
            },
          }),
        ],
      },
      {
        path: '/monitor',
        name: 'BootstrapGroupMonitor',
        meta: {
          orderNo: 5,
        },
        children: [
          asRouteRecordRaw({
            path: 'overview',
            name: 'MonitorOverviewIndex',
            meta: {
              orderNo: 5,
              tabTitle: {
                'zh-CN': '监控中心 - 服务概览',
                'en-US': 'Monitor - Overview',
              },
              titleKey: 'monitor.route.overview.title',
            },
          }),
        ],
      },
      {
        path: '/notifications',
        name: 'NotificationList',
        meta: {
          hiddenMenu: true,
          tabTitle: {
            'zh-CN': '通知中心',
            'en-US': 'Notifications',
          },
          titleKey: 'notification.route.list.title',
        },
      },
    ] as RouteRecordRaw[];

    expect(buildDashboardQuickActionLinks(routes, 'zh-CN')).toEqual([
      {
        icon: undefined,
        id: 'MonitorOverviewIndex',
        module_key: 'monitor',
        order: 5,
        route_location: '/monitor/overview',
        title: '监控中心 - 服务概览',
        title_key: 'monitor.route.overview.title',
      },
      {
        icon: 'layers',
        id: 'ContainerList',
        module_key: 'ops',
        order: 10,
        route_location: '/ops/containers',
        title: '容器管理',
        title_key: 'container.route.list.title',
      },
      {
        icon: 'secured',
        id: 'RoleListIndex',
        module_key: 'rbac',
        order: 20,
        route_location: '/access-control/roles',
        title: '访问控制 - 角色管理',
        title_key: 'rbac.role.list.title',
      },
    ]);
  });

  it('uses the requested locale instead of hard-coding english titles', () => {
    const routes = [
      {
        path: '/monitor',
        name: 'BootstrapGroupMonitor',
        children: [
          asRouteRecordRaw({
            path: 'overview',
            name: 'MonitorOverviewIndex',
            meta: {
              tabTitle: {
                'zh-CN': '监控中心 - 服务概览',
                'en-US': 'Monitor - Overview',
              },
              titleKey: 'monitor.route.overview.title',
            },
          }),
        ],
      },
    ] as RouteRecordRaw[];

    expect(buildDashboardQuickActionLinks(routes, 'en-US')).toEqual([
      {
        icon: undefined,
        id: 'MonitorOverviewIndex',
        module_key: 'monitor',
        order: 0,
        route_location: '/monitor/overview',
        title: 'Monitor - Overview',
        title_key: 'monitor.route.overview.title',
      },
    ]);
  });
});
