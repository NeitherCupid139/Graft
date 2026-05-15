<template>
  <div class="user-page">
    <t-row :gutter="[24, 24]">
      <t-col :span="12">
        <t-card class="summary-card" :bordered="false" title="用户列表">
          <div class="summary-metric">
            <span class="summary-metric__label">当前返回用户数</span>
            <span class="summary-metric__value">{{ users.length }}</span>
          </div>
          <div class="summary-hint">当前页面直接消费 `GET /api/users`，不再保留 starter demo 个人中心数据。</div>
        </t-card>
      </t-col>
      <t-col :span="12">
        <t-card class="summary-card" :bordered="false" title="接口说明">
          <div class="summary-meta">
            <span>接口路径：`/api/users`</span>
            <span>字段：`id / username / display / created_at / updated_at`</span>
          </div>
          <t-button theme="primary" variant="outline" :loading="loading" @click="fetchUsers">刷新列表</t-button>
        </t-card>
      </t-col>
    </t-row>

    <t-card class="table-card" :bordered="false" title="用户数据">
      <t-table row-key="id" :data="users" :columns="columns" :loading="loading" size="medium">
        <template #empty>
          <t-empty description="暂无用户数据" />
        </template>
      </t-table>
    </t-card>
  </div>
</template>
<script setup lang="ts">
import { MessagePlugin, type TdBaseTableProps } from 'tdesign-vue-next';
import { computed, onMounted, ref } from 'vue';

import type { UserListItem } from '@/api/model/userModel';
import { getUsers } from '@/api/user';

defineOptions({
  name: 'UsersIndex',
});

const users = ref<UserListItem[]>([]);
const loading = ref(false);

const columns = computed<TdBaseTableProps['columns']>(() => [
  {
    title: 'ID',
    colKey: 'id',
    width: 100,
  },
  {
    title: '用户名',
    colKey: 'username',
    minWidth: 180,
  },
  {
    title: '显示名',
    colKey: 'display',
    minWidth: 180,
  },
  {
    title: '创建时间',
    colKey: 'created_at',
    minWidth: 220,
  },
  {
    title: '更新时间',
    colKey: 'updated_at',
    minWidth: 220,
  },
]);

async function fetchUsers() {
  loading.value = true;
  try {
    // 后端契约假设：GET /api/users 始终返回 { items: UserListItem[] }，且时间字段已可直接展示。
    const response = await getUsers();
    users.value = response.items;
  } catch (error) {
    users.value = [];
    MessagePlugin.error(error instanceof Error ? error.message : '用户列表加载失败');
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchUsers();
});
</script>
<style lang="less" scoped>
@import './index.less';
</style>
