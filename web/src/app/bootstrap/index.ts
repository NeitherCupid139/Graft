/* eslint-disable simple-import-sort/imports */
import { createApp } from 'vue';
import TDesign from 'tdesign-vue-next';

import App from '@/App.vue';
import { i18n } from '@/locales';
import router from '@/router';
import { store, useTabsRouterStore } from '@/store';

import { registerPermissionDirective } from './permission-directive';
import { registerRouteGuards } from './route-guards';

import 'tdesign-vue-next/es/style/index.css';
import '@/style/index.less';

// bootstrapApp owns the single startup path for the real web runtime.
export function bootstrapApp() {
  registerRouteGuards(router);

  const app = createApp(App);
  const tabsRouterStore = useTabsRouterStore(store);

  // 自愈上一次异常中断时遗留的 tabs 刷新态，避免内容区被全局 loading 永久覆盖。
  tabsRouterStore.healPersistedState();

  app.use(TDesign);
  app.use(store);
  app.use(router);
  app.use(i18n);
  // 权限指令只消费 bootstrap 权限快照，不引入第二套前端鉴权真值。
  registerPermissionDirective(app);

  app.mount('#app');

  return app;
}
