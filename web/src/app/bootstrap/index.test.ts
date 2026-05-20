import { beforeEach, describe, expect, it, vi } from 'vitest';

const healPersistedState = vi.fn();
const registerRouteGuards = vi.fn();
const registerPermissionDirective = vi.fn();
const useMock = vi.fn();
const mountMock = vi.fn();

vi.mock('vue', () => ({
  createApp: () => ({
    use: useMock,
    mount: mountMock,
  }),
}));

vi.mock('tdesign-vue-next', () => ({
  default: {},
}));

vi.mock('@/App.vue', () => ({
  default: {},
}));

vi.mock('@/router', () => ({
  default: {},
}));

vi.mock('@/locales', () => ({
  i18n: {},
}));

vi.mock('@/store', () => ({
  store: {},
  useTabsRouterStore: () => ({
    healPersistedState,
  }),
}));

vi.mock('./route-guards', () => ({
  registerRouteGuards,
}));

vi.mock('./permission-directive', () => ({
  registerPermissionDirective,
}));

describe('bootstrapApp', () => {
  beforeEach(() => {
    vi.resetModules();
    healPersistedState.mockReset();
    registerRouteGuards.mockReset();
    registerPermissionDirective.mockReset();
    useMock.mockReset();
    mountMock.mockReset();
  });

  it('heals persisted tab refresh residue before mounting the app', async () => {
    const { bootstrapApp } = await import('./index');

    bootstrapApp();

    expect(registerRouteGuards).toHaveBeenCalledTimes(1);
    expect(healPersistedState).toHaveBeenCalledTimes(1);
    expect(useMock).toHaveBeenCalledTimes(4);
    expect(registerPermissionDirective).toHaveBeenCalledTimes(1);
    expect(mountMock).toHaveBeenCalledWith('#app');
  });
});
