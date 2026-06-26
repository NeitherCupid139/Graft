<template>
  <div class="shell-header-resource-actions">
    <t-tooltip placement="bottom" :content="t('layout.header.code')">
      <t-button theme="default" shape="square" variant="text" @click="openRepository">
        <t-icon name="logo-github" />
      </t-button>
    </t-tooltip>
    <t-tooltip placement="bottom" :content="t('layout.header.apiDocs')">
      <t-button theme="default" shape="square" variant="text" @click="openApiDocs">
        <t-icon name="book-open" />
      </t-button>
    </t-tooltip>
    <t-tooltip placement="bottom" :content="t('layout.header.help')">
      <t-button theme="default" shape="square" variant="text" @click="openHelp">
        <t-icon name="help-circle" />
      </t-button>
    </t-tooltip>
  </div>
</template>
<script setup lang="ts">
import { t } from '@/locales';

const REPOSITORY_URL = 'https://github.com/GeWuYou/Graft';

const openRepository = () => {
  window.open(REPOSITORY_URL);
};

const openApiDocs = () => {
  window.open('/docs');
};

const openHelp = () => {
  window.open(`${REPOSITORY_URL}#readme`);
};
</script>
<style lang="less" scoped>
.shell-header-resource-actions {
  align-items: center;
  display: inline-flex;
  gap: var(--graft-density-gap-4);
}
</style>