import { mount } from '@vue/test-utils';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { markRaw } from 'vue';

import Content from './Content.vue';

const routeState = vi.hoisted(() => ({
  meta: {},
  path: '/access-control/roles',
  fullPath: '/access-control/roles',
}));
const tabStoreState = vi.hoisted(() => ({
  activeTabKey: '/access-control/roles',
  refreshing: false,
  tabRouters: [
    {
      isAlive: true,
      meta: {},
      name: 'RoleListIndex',
    },
  ],
}));
const tabStoreProxy = vi.hoisted(() => ({
  value: null as null | typeof tabStoreState,
}));

vi.mock('vue-router', () => ({
  useRoute: () => routeState,
}));

vi.mock('@/store', async () => {
  const { reactive } = await import('vue');
  tabStoreProxy.value = reactive(tabStoreState);

  return {
    useTabsRouterStore: () => tabStoreProxy.value,
  };
});

describe('Content', () => {
  beforeEach(() => {
    routeState.path = '/access-control/roles';
    routeState.fullPath = '/access-control/roles';
    routeState.meta = {};
    tabStoreState.activeTabKey = '/access-control/roles';
    tabStoreState.refreshing = false;
  });

  it('keys rendered route content by the active tab key', async () => {
    const wrapper = mount(Content, {
      global: {
        stubs: {
          RouterView: {
            template: '<slot :Component="Component" />',
            data() {
              return {
                Component: markRaw({
                  name: 'RouteContentProbe',
                  template: '<div data-testid="route-content">content</div>',
                }),
              };
            },
          },
          Transition: {
            template: '<slot />',
          },
          KeepAlive: {
            props: ['include'],
            template: '<slot />',
          },
          FramePage: true,
          TLoading: true,
        },
      },
    });

    expect(wrapper.findComponent({ name: 'RouteContentProbe' }).vm.$.vnode.key).toBe('/access-control/roles');

    tabStoreProxy.value!.activeTabKey = '/access-control/roles#copy-1';
    await wrapper.vm.$nextTick();

    expect(wrapper.findComponent({ name: 'RouteContentProbe' }).vm.$.vnode.key).toBe('/access-control/roles#copy-1');
  });
});
