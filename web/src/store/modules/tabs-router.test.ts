import { createPinia, setActivePinia } from 'pinia';
import { beforeEach, describe, expect, it } from 'vitest';

import { useTabsRouterStore } from './tabs-router';

describe('useTabsRouterStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('uses the neutral root entry as the preserved home tab', () => {
    const tabsRouterStore = useTabsRouterStore();

    expect(tabsRouterStore.tabRouters).toHaveLength(1);
    expect(tabsRouterStore.tabRouters[0]?.path).toBe('/');
    expect(tabsRouterStore.tabRouters[0]?.name).toBe('RootEntry');
  });

  it('keeps refresh state ephemeral and restores the tab after refresh completes', () => {
    const tabsRouterStore = useTabsRouterStore();

    tabsRouterStore.appendTabRouterList({
      path: '/users',
      name: 'UserList',
      meta: {
        keepAlive: true,
      },
    });

    tabsRouterStore.startTabRefresh(1);
    expect(tabsRouterStore.refreshing).toBe(true);
    expect(tabsRouterStore.tabRouters[1]?.isAlive).toBe(false);

    tabsRouterStore.finishTabRefresh(1);
    expect(tabsRouterStore.refreshing).toBe(false);
    expect(tabsRouterStore.tabRouters[1]?.isAlive).toBe(true);
  });

  it('heals persisted refresh residue on startup', () => {
    const tabsRouterStore = useTabsRouterStore();

    tabsRouterStore.appendTabRouterList({
      path: '/users',
      name: 'UserList',
      meta: {
        keepAlive: true,
      },
    });

    tabsRouterStore.startTabRefresh(1);
    expect(tabsRouterStore.refreshing).toBe(true);
    expect(tabsRouterStore.tabRouters[1]?.isAlive).toBe(false);

    tabsRouterStore.healPersistedState();
    expect(tabsRouterStore.refreshing).toBe(false);
    expect(tabsRouterStore.tabRouters[0]?.isAlive).toBe(true);
    expect(tabsRouterStore.tabRouters[1]?.isAlive).toBe(true);
  });
});
