<template>
  <div class="permission-page" data-page-type="list-form-detail">
    <management-page-header :title="t('rbac.permissionList.listTitle')" :description="t('rbac.permissionList.hint')">
      <template #eyebrow>{{ t('menu.access_control.title') }}</template>
      <template #actions>
        <t-button theme="default" variant="outline" :loading="loading" @click="fetchPermissions">
          {{ t('rbac.permissionList.refresh') }}
        </t-button>
      </template>
    </management-page-header>

    <management-toolbar>
      <template #filters>
        <t-input
          v-model="filters.keyword"
          clearable
          class="toolbar__search"
          :placeholder="t('rbac.permissionList.toolbar.searchPlaceholder')"
        />
        <t-select
          v-model="filters.category"
          clearable
          class="toolbar__select"
          :options="categoryOptions"
          :placeholder="t('rbac.permissionList.toolbar.modulePlaceholder')"
        />
        <t-button variant="text" @click="resetFilters">
          {{ t('rbac.permissionList.toolbar.clearFilters') }}
        </t-button>
      </template>
      <template #actions>
        <t-button theme="default" variant="outline" @click="columnDrawerVisible = true">
          {{ t('rbac.permissionList.columnSettings') }}
        </t-button>
        <t-button theme="default" variant="outline" :loading="loading" @click="fetchPermissions">
          {{ t('rbac.permissionList.refresh') }}
        </t-button>
      </template>
    </management-toolbar>

    <management-table-card>
      <template #head>
        <div class="table-head">
          <div>
            <p class="table-head__summary">
              {{ t('rbac.permissionList.summary', { count: filteredPermissions.length }) }}
            </p>
            <p class="table-head__description">{{ t('rbac.permissionList.readonlyNotice') }}</p>
          </div>
        </div>
      </template>

      <management-empty-state
        v-if="listError && !loading"
        tone="error"
        :title="t('rbac.permissionList.errorTitle')"
        :description="listError"
      >
        <template #actions>
          <t-button theme="primary" variant="outline" @click="fetchPermissions">
            {{ t('rbac.permissionList.retry') }}
          </t-button>
        </template>
      </management-empty-state>

      <t-table
        v-else
        row-key="id"
        :data="filteredPermissions"
        :columns="visibleColumns"
        :loading="loading"
        cell-empty-content="-"
      >
        <template #permission="{ row }">
          <div class="permission-cell">
            <span class="permission-cell__name">{{ row.display }}</span>
            <span class="permission-cell__code">{{ row.code }}</span>
          </div>
        </template>

        <template #category="{ row }">
          <t-tag theme="default" variant="light">{{ row.category || '-' }}</t-tag>
        </template>

        <template #role_count="{ row }">
          <span>{{ roleUsageMap[row.id] ?? '-' }}</span>
        </template>

        <template #operation="{ row }">
          <div class="table-actions">
            <t-button size="small" variant="outline" @click="showReadonlyMessage(row.code)">
              {{ t('rbac.permissionList.detail') }}
            </t-button>
          </div>
        </template>

        <template #empty>
          <management-empty-state
            :title="t('rbac.permissionList.emptyTitle')"
            :description="t('rbac.permissionList.empty')"
          />
        </template>
      </t-table>
    </management-table-card>

    <t-drawer
      v-model:visible="columnDrawerVisible"
      :header="t('rbac.permissionList.columnSettings')"
      size="360px"
      placement="right"
      destroy-on-close
    >
      <div class="drawer-panel">
        <t-checkbox-group v-model="visibleColumnKeys">
          <div class="column-grid">
            <t-checkbox v-for="column in columnSettingOptions" :key="column.value" :value="column.value">
              {{ column.label }}
            </t-checkbox>
          </div>
        </t-checkbox-group>
      </div>
    </t-drawer>
  </div>
</template>
<script setup lang="ts">
import { MessagePlugin, type TdBaseTableProps } from 'tdesign-vue-next';
import { computed, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import {
  ManagementEmptyState,
  ManagementPageHeader,
  ManagementTableCard,
  ManagementToolbar,
} from '@/shared/components/management';

import { getPermissions, getRoles } from '../../api/rbac';
import type { PermissionListItem } from '../../types/rbac';

defineOptions({
  name: 'PermissionIndex',
});

type PermissionFilters = {
  keyword: string;
  category: string;
};

const { t, locale } = useI18n();
const loading = ref(false);
const listError = ref('');
const permissions = ref<PermissionListItem[]>([]);
const roleUsageMap = ref<Record<number, number>>({});
const filters = ref<PermissionFilters>({
  keyword: '',
  category: '',
});
const columnDrawerVisible = ref(false);
const visibleColumnKeys = ref(['permission', 'category', 'code', 'role_count', 'operation']);

const categoryOptions = computed(() => {
  const categories = Array.from(new Set(permissions.value.map((item) => item.category).filter(Boolean))).sort();
  return categories.map((category) => ({ label: category, value: category }));
});

const columnSettingOptions = computed(() => [
  { label: t('rbac.permissionList.columns.permission'), value: 'permission' },
  { label: t('rbac.permissionList.columns.module'), value: 'category' },
  { label: t('rbac.permissionList.columns.code'), value: 'code' },
  { label: t('rbac.permissionList.columns.roleCount'), value: 'role_count' },
  { label: t('components.commonTable.operation'), value: 'operation' },
]);

const filteredPermissions = computed(() => {
  const keyword = filters.value.keyword.trim().toLowerCase();

  return permissions.value.filter((item) => {
    if (filters.value.category && item.category !== filters.value.category) {
      return false;
    }

    if (!keyword) {
      return true;
    }

    return `${item.code} ${item.display} ${item.description ?? ''} ${item.category}`.toLowerCase().includes(keyword);
  });
});

const visibleColumns = computed<TdBaseTableProps['columns']>(() => {
  void locale.value;

  const allColumns: TdBaseTableProps['columns'] = [
    { title: t('rbac.permissionList.columns.permission'), colKey: 'permission', minWidth: 260, fixed: 'left' },
    { title: t('rbac.permissionList.columns.module'), colKey: 'category', width: 180 },
    { title: t('rbac.permissionList.columns.code'), colKey: 'code', minWidth: 240 },
    { title: t('rbac.permissionList.columns.roleCount'), colKey: 'role_count', width: 140 },
    { title: t('components.commonTable.operation'), colKey: 'operation', width: 140, fixed: 'right' },
  ];

  const visibleKeys = new Set(visibleColumnKeys.value);
  return allColumns.filter((column) => visibleKeys.has(String(column.colKey)));
});

async function fetchPermissions() {
  loading.value = true;
  listError.value = '';

  try {
    const [permissionResult, roleResult] = await Promise.all([getPermissions(), getRoles()]);
    permissions.value = permissionResult.items;

    const usageMap: Record<number, number> = {};
    roleResult.items.forEach((role) => {
      if (typeof role.permission_count !== 'number') {
        return;
      }
      permissions.value.forEach((permission) => {
        usageMap[permission.id] = usageMap[permission.id] ?? 0;
      });
    });
    roleUsageMap.value = usageMap;
  } catch (error) {
    permissions.value = [];
    roleUsageMap.value = {};
    listError.value = error instanceof Error ? error.message : t('rbac.permissionList.loadFailed');
    MessagePlugin.error(listError.value);
  } finally {
    loading.value = false;
  }
}

function resetFilters() {
  filters.value = {
    keyword: '',
    category: '',
  };
}

function showReadonlyMessage(permissionCode: string) {
  MessagePlugin.warning(`${t('rbac.permissionList.readonlyHint')}: ${permissionCode}`);
}

onMounted(() => {
  fetchPermissions();
});
</script>
<style scoped lang="less">
.permission-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.toolbar__search {
  width: min(100%, 320px);
}

.toolbar__select {
  width: min(100%, 220px);
}

.table-head,
.table-actions {
  align-items: center;
  display: flex;
  gap: 12px;
  justify-content: space-between;
}

.table-head__summary,
.table-head__description,
.permission-cell__code {
  color: var(--td-text-color-secondary);
  margin: 0;
}

.permission-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.permission-cell__name {
  color: var(--td-text-color-primary);
  font: var(--td-font-title-small);
}

.drawer-panel,
.column-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

@media (width <= 768px) {
  .toolbar__search,
  .toolbar__select {
    width: 100%;
  }

  .table-head,
  .table-actions {
    align-items: stretch;
    flex-direction: column;
  }
}
</style>
