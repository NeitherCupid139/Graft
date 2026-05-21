<template>
  <div class="user-page" data-page-type="list-form-detail">
    <header class="admin-page-header">
      <div class="admin-page-header__copy">
        <h1 class="admin-page-header__title">{{ t('user.userList.listTitle') }}</h1>
        <p class="admin-page-header__description">{{ t('user.userList.hint') }}</p>
      </div>
    </header>

    <section class="admin-surface">
      <div class="toolbar">
        <div class="toolbar__filters">
          <t-input
            v-model="filters.keyword"
            clearable
            class="toolbar__search"
            :placeholder="t('user.userList.toolbar.searchPlaceholder')"
          />
          <t-select
            v-model="filters.status"
            clearable
            class="toolbar__select"
            :options="statusOptions"
            :placeholder="t('user.userList.toolbar.statusPlaceholder')"
          />
          <t-select
            v-model="filters.roleId"
            clearable
            class="toolbar__select"
            :options="roleOptions"
            :loading="roleCatalogLoading"
            :placeholder="t('user.userList.toolbar.rolePlaceholder')"
          />
          <t-button variant="text" @click="resetFilters">
            {{ t('user.userList.toolbar.clearFilters') }}
          </t-button>
        </div>
        <div class="toolbar__actions">
          <t-button theme="default" variant="outline" :loading="loading" data-testid="user-refresh" @click="fetchUsers">
            {{ t('user.userList.refresh') }}
          </t-button>
          <t-button theme="default" variant="outline" @click="columnDrawerVisible = true">
            {{ t('user.userList.columnSettings') }}
          </t-button>
          <t-button theme="primary" data-testid="user-create" @click="handleUnavailableAction('create')">
            {{ t('user.userList.create') }}
          </t-button>
        </div>
      </div>
    </section>

    <section class="admin-surface admin-surface--table">
      <div class="table-head">
        <p class="table-head__summary">{{ t('user.userList.summary', { count: filteredUsers.length }) }}</p>
        <t-button v-if="hasActiveFilters" variant="text" @click="resetFilters">
          {{ t('user.userList.toolbar.clearFilters') }}
        </t-button>
      </div>

      <div v-if="listError && !loading" class="state-panel state-panel--error" data-testid="user-list-error">
        <p class="state-panel__title">{{ t('user.userList.errorTitle') }}</p>
        <p class="state-panel__description">{{ listError }}</p>
        <t-button theme="primary" variant="outline" @click="fetchUsers">
          {{ t('user.userList.retry') }}
        </t-button>
      </div>

      <t-table
        v-else
        row-key="id"
        :data="filteredUsers"
        :columns="visibleColumns"
        :loading="loading"
        cell-empty-content="-"
      >
        <template #user="{ row }">
          <div class="user-cell">
            <div class="user-cell__avatar">{{ userInitial(row.display || row.username) }}</div>
            <div class="user-cell__meta">
              <span class="user-cell__display">{{ row.display || row.username }}</span>
              <span class="user-cell__username">@{{ row.username }}</span>
            </div>
          </div>
        </template>

        <template #status="{ row }">
          <t-tag :theme="statusTheme(row.status)" variant="light">
            {{ statusLabel(row.status) }}
          </t-tag>
        </template>

        <template #roles="{ row }">
          <div class="role-tag-list">
            <template v-if="roleSummaryLoading[row.id]">
              <t-tag theme="default" variant="light">{{ t('user.userList.roleSummary.loading') }}</t-tag>
            </template>
            <template v-else-if="roleSummaryErrors[row.id]">
              <span class="table-muted">{{ t('user.userList.roleSummary.unavailable') }}</span>
            </template>
            <template v-else-if="resolveUserRoles(row.id).length > 0">
              <t-tag
                v-for="role in resolveUserRoles(row.id)"
                :key="role.id"
                theme="default"
                variant="light-outline"
                size="small"
              >
                {{ role.display }}
              </t-tag>
            </template>
            <span v-else class="table-muted">{{ t('user.userList.roleSummary.empty') }}</span>
          </div>
        </template>

        <template #created_at="{ row }">
          <span>{{ formatTimestamp(row.created_at) }}</span>
        </template>

        <template #updated_at="{ row }">
          <span>{{ formatTimestamp(row.updated_at) }}</span>
        </template>

        <template #operation="{ row }">
          <div class="table-actions">
            <t-button
              v-if="canReadUserRoles"
              size="small"
              theme="primary"
              variant="outline"
              data-testid="user-manage-roles"
              @click="handleOpenUserRoleDrawer(row)"
            >
              {{ t('user.userList.assignRoles') }}
            </t-button>
            <t-button size="small" theme="default" variant="outline" @click="handleUnavailableAction('edit')">
              {{ t('user.userList.edit') }}
            </t-button>
            <t-button size="small" theme="default" variant="outline" @click="handleUnavailableAction('more')">
              {{ t('user.userList.more') }}
            </t-button>
          </div>
        </template>

        <template #empty>
          <t-empty :description="t('user.userList.empty')" />
        </template>
      </t-table>
    </section>

    <t-drawer
      v-model:visible="userRoleDrawerVisible"
      :header="t('user.userList.roleDialog.title')"
      size="520px"
      placement="right"
      destroy-on-close
    >
      <div class="drawer-panel" data-testid="user-role-drawer">
        <div class="drawer-summary">
          <div class="user-cell user-cell--drawer">
            <div class="user-cell__avatar">{{ userInitial(selectedUser?.display || selectedUser?.username) }}</div>
            <div class="user-cell__meta">
              <span class="user-cell__display">{{ selectedUser?.display || '-' }}</span>
              <span class="user-cell__username">@{{ selectedUser?.username || '-' }}</span>
            </div>
          </div>
          <div class="drawer-summary__grid">
            <div class="drawer-summary__item">
              <span class="drawer-summary__label">{{ t('user.userList.columns.status') }}</span>
              <t-tag :theme="statusTheme(selectedUser?.status)" variant="light">
                {{ statusLabel(selectedUser?.status) }}
              </t-tag>
            </div>
            <div class="drawer-summary__item">
              <span class="drawer-summary__label">{{ t('user.userList.columns.createdAt') }}</span>
              <span>{{ formatTimestamp(selectedUser?.created_at) }}</span>
            </div>
            <div class="drawer-summary__item">
              <span class="drawer-summary__label">{{ t('user.userList.columns.updatedAt') }}</span>
              <span>{{ formatTimestamp(selectedUser?.updated_at) }}</span>
            </div>
          </div>
        </div>

        <section class="drawer-section">
          <div class="drawer-section__head">
            <h3>{{ t('user.userList.roleDialog.currentRolesTitle') }}</h3>
            <span class="table-muted">{{ t('user.userList.roleDialog.roleSummary', { count: roles.length }) }}</span>
          </div>
          <div class="role-tag-list">
            <template v-if="currentUserRoles.length > 0">
              <t-tag v-for="role in currentUserRoles" :key="role.id" theme="default" variant="light-outline">
                {{ role.display }}
              </t-tag>
            </template>
            <span v-else class="table-muted">{{ t('user.userList.roleDialog.noAssignedRoles') }}</span>
          </div>
        </section>

        <div v-if="roleLoadWarning" class="inline-warning">
          <span>{{ roleLoadWarning }}</span>
          <t-button variant="text" theme="primary" :loading="loadingRoleDialogData" @click="retryUserRoleDrawerLoad">
            {{ t('user.userList.roleDialog.retry') }}
          </t-button>
        </div>

        <section class="drawer-section">
          <div class="drawer-section__head">
            <h3>{{ t('user.userList.roleDialog.availableRolesTitle') }}</h3>
          </div>
          <t-checkbox-group
            v-model="selectedRoleIds"
            :disabled="loadingRoleDialogData || !roleSelectionReady || !canAssignUserRoles"
            data-testid="role-checkbox-group"
          >
            <div class="selection-grid">
              <label v-for="role in roles" :key="role.id" class="selection-card">
                <t-checkbox :value="role.id">
                  <div class="selection-card__body">
                    <div class="selection-card__head">
                      <span class="selection-card__title">{{ role.display }}</span>
                      <t-tag :theme="role.builtin ? 'warning' : 'default'" variant="light" size="small">
                        {{
                          role.builtin
                            ? t('user.userList.roleDialog.builtinYes')
                            : t('user.userList.roleDialog.builtinNo')
                        }}
                      </t-tag>
                    </div>
                    <span class="selection-card__code">{{ role.name }}</span>
                    <span class="selection-card__description">
                      {{ role.description || t('user.userList.roleDialog.emptyDescription') }}
                    </span>
                  </div>
                </t-checkbox>
              </label>
            </div>
          </t-checkbox-group>
          <t-empty
            v-if="roles.length === 0 && !roleCatalogLoading"
            :description="t('user.userList.roleDialog.empty')"
          />
        </section>

        <div class="drawer-actions">
          <t-button variant="outline" data-testid="user-role-cancel" @click="closeUserRoleDrawer">
            {{ t('user.userList.roleDialog.cancel') }}
          </t-button>
          <t-button
            theme="primary"
            data-testid="user-role-save"
            :disabled="!canSubmitRoleAssignment"
            :loading="submittingRoles"
            @click="submitUserRoleAssignment"
          >
            {{ t('user.userList.roleDialog.confirm') }}
          </t-button>
        </div>
      </div>
    </t-drawer>

    <t-drawer
      v-model:visible="columnDrawerVisible"
      :header="t('user.userList.columnSettings')"
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

import { RBAC_PERMISSION_CODE } from '@/modules/rbac/contract/permissions';
import type { RoleListItem } from '@/modules/rbac/contract/role';
import { usePermissionStore } from '@/store';

import { assignUserRoles, getRoles, getUserRoleBindings } from '../api/user-roles';
import { getUsers } from '../api/users';
import { USER_PERMISSION_CODE } from '../contract/permissions';
import type { UserStatus } from '../contract/status';
import { USER_STATUS } from '../contract/status';
import type { UserListItem } from '../types/user';

defineOptions({
  name: 'UsersIndex',
});

type UserFilters = {
  keyword: string;
  roleId: number | undefined;
  status: '' | UserStatus;
};

type UserRow = UserListItem;

const DEFAULT_VISIBLE_COLUMNS = ['user', 'status', 'roles', 'created_at', 'updated_at', 'operation'];

const { t, locale } = useI18n();
const permissionStore = usePermissionStore();
const users = ref<UserRow[]>([]);
const roles = ref<RoleListItem[]>([]);
const loading = ref(false);
const listError = ref('');
const roleCatalogLoading = ref(false);
const roleCatalogError = ref('');
const roleSummaryRequestId = ref(0);
const roleBindings = ref<Record<number, number[]>>({});
const roleSummaryLoading = ref<Record<number, boolean>>({});
const roleSummaryErrors = ref<Record<number, boolean>>({});
const filters = ref<UserFilters>({
  keyword: '',
  roleId: undefined,
  status: '',
});
const visibleColumnKeys = ref<string[]>([...DEFAULT_VISIBLE_COLUMNS]);
const columnDrawerVisible = ref(false);
const userRoleDrawerVisible = ref(false);
const selectedUser = ref<UserRow | null>(null);
const selectedRoleIds = ref<number[]>([]);
const loadingRoleSelection = ref(false);
const submittingRoles = ref(false);
const roleSelectionReady = ref(false);
const roleLoadWarning = ref('');
const drawerSession = ref(0);

const userPermissionCodes = USER_PERMISSION_CODE;
const rbacPermissionCodes = RBAC_PERMISSION_CODE;
const canReadUserRoles = computed(() => permissionStore.hasPermission(rbacPermissionCodes.USER_ROLE_READ));
const canAssignUserRoles = computed(() => permissionStore.hasPermission(rbacPermissionCodes.USER_ROLE_ASSIGN));
const canShowOperationColumn = computed(() =>
  permissionStore.hasAnyPermission([
    userPermissionCodes.UPDATE,
    userPermissionCodes.CREATE,
    rbacPermissionCodes.USER_ROLE_READ,
  ]),
);
const loadingRoleDialogData = computed(() => roleCatalogLoading.value || loadingRoleSelection.value);
const canSubmitRoleAssignment = computed(
  () => canAssignUserRoles.value && roleSelectionReady.value && selectedUser.value !== null,
);
const hasActiveFilters = computed(
  () => Boolean(filters.value.keyword.trim()) || Boolean(filters.value.status) || filters.value.roleId !== undefined,
);

const statusOptions = computed(() => [
  { label: t('user.userList.toolbar.statusAll'), value: '' },
  { label: t('user.userList.status.enabled'), value: USER_STATUS.ENABLED },
  { label: t('user.userList.status.disabled'), value: USER_STATUS.DISABLED },
]);

const roleOptions = computed(() =>
  roles.value.map((role) => ({
    label: role.display,
    value: role.id,
  })),
);

const columnSettingOptions = computed(() => [
  { label: t('user.userList.columns.user'), value: 'user' },
  { label: t('user.userList.columns.status'), value: 'status' },
  { label: t('user.userList.columns.roles'), value: 'roles' },
  { label: t('user.userList.columns.createdAt'), value: 'created_at' },
  { label: t('user.userList.columns.updatedAt'), value: 'updated_at' },
  { label: t('components.commonTable.operation'), value: 'operation' },
]);

const filteredUsers = computed(() => {
  const keyword = filters.value.keyword.trim().toLowerCase();

  return users.value.filter((user) => {
    if (keyword) {
      const haystack = `${user.username} ${user.display}`.toLowerCase();
      if (!haystack.includes(keyword)) {
        return false;
      }
    }

    if (filters.value.status && normalizeUserStatus(user.status) !== filters.value.status) {
      return false;
    }

    if (filters.value.roleId !== undefined) {
      const assignedRoleIds = roleBindings.value[user.id] ?? [];
      if (!assignedRoleIds.includes(filters.value.roleId)) {
        return false;
      }
    }

    return true;
  });
});

const currentUserRoles = computed(() => {
  const roleIDs = selectedRoleIds.value;
  return roles.value.filter((role) => roleIDs.includes(role.id));
});

const columns = computed<TdBaseTableProps['columns']>(() => {
  void locale.value;

  const allColumns: TdBaseTableProps['columns'] = [
    {
      title: t('user.userList.columns.user'),
      colKey: 'user',
      minWidth: 240,
    },
    {
      title: t('user.userList.columns.status'),
      colKey: 'status',
      width: 120,
    },
    {
      title: t('user.userList.columns.roles'),
      colKey: 'roles',
      minWidth: 220,
    },
    {
      title: t('user.userList.columns.createdAt'),
      colKey: 'created_at',
      width: 200,
    },
    {
      title: t('user.userList.columns.updatedAt'),
      colKey: 'updated_at',
      width: 200,
    },
  ];

  if (canShowOperationColumn.value) {
    allColumns.push({
      title: t('components.commonTable.operation'),
      colKey: 'operation',
      width: 260,
      fixed: 'right',
    });
  }

  const visibleKeys = new Set(visibleColumnKeys.value);
  return allColumns.filter((column) => visibleKeys.has(String(column.colKey)));
});

const visibleColumns = computed(() => {
  if (canShowOperationColumn.value) {
    return columns.value;
  }

  return (columns.value ?? []).filter((column) => column?.colKey !== 'operation');
});

async function fetchUsers() {
  loading.value = true;
  listError.value = '';
  roleSummaryErrors.value = {};
  roleSummaryLoading.value = {};

  try {
    const response = await getUsers();
    users.value = response.items;

    if (canReadUserRoles.value) {
      void hydrateUserRoleSummaries(response.items);
    } else {
      roleBindings.value = {};
    }
  } catch (error) {
    users.value = [];
    listError.value = error instanceof Error ? error.message : t('user.userList.loadFailed');
    MessagePlugin.error(listError.value);
  } finally {
    loading.value = false;
  }
}

async function loadRoleCatalog() {
  roleCatalogLoading.value = true;

  try {
    const response = await getRoles();
    roles.value = response.items;
    roleCatalogError.value = '';
  } catch (error) {
    roles.value = [];
    roleCatalogError.value = error instanceof Error ? error.message : t('user.userList.roleSummary.loadFailed');
    throw error;
  } finally {
    roleCatalogLoading.value = false;
  }
}

async function hydrateUserRoleSummaries(userItems: UserRow[]) {
  const requestId = roleSummaryRequestId.value + 1;
  roleSummaryRequestId.value = requestId;

  try {
    await loadRoleCatalog();
  } catch {
    return;
  }

  const nextLoading = Object.fromEntries(userItems.map((user) => [user.id, true]));
  roleSummaryLoading.value = nextLoading;
  roleBindings.value = {};
  roleSummaryErrors.value = {};

  const results = await Promise.allSettled(userItems.map((user) => getUserRoleBindings(user.id)));
  if (roleSummaryRequestId.value !== requestId) {
    return;
  }

  const nextBindings: Record<number, number[]> = {};
  const nextErrors: Record<number, boolean> = {};
  const nextLoadingDone: Record<number, boolean> = {};

  userItems.forEach((user, index) => {
    const result = results[index];
    nextLoadingDone[user.id] = false;

    if (result?.status === 'fulfilled') {
      nextBindings[user.id] = result.value.role_ids;
      return;
    }

    nextErrors[user.id] = true;
  });

  roleBindings.value = nextBindings;
  roleSummaryErrors.value = nextErrors;
  roleSummaryLoading.value = nextLoadingDone;
}

function resolveUserRoles(userId: number) {
  const assignedRoleIds = new Set(roleBindings.value[userId] ?? []);
  return roles.value.filter((role) => assignedRoleIds.has(role.id));
}

function resetFilters() {
  filters.value = {
    keyword: '',
    roleId: undefined,
    status: '',
  };
}

function formatTimestamp(value?: string | null) {
  if (!value) {
    return '-';
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return new Intl.DateTimeFormat(locale.value, {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(date);
}

function userInitial(value?: string | null) {
  if (!value) {
    return '?';
  }

  return value.trim().slice(0, 1).toUpperCase();
}

function normalizeUserStatus(status?: string | null) {
  return status === USER_STATUS.DISABLED ? USER_STATUS.DISABLED : USER_STATUS.ENABLED;
}

function statusLabel(status?: string | null) {
  return normalizeUserStatus(status) === USER_STATUS.DISABLED
    ? t('user.userList.status.disabled')
    : t('user.userList.status.enabled');
}

function statusTheme(status?: string | null) {
  return normalizeUserStatus(status) === USER_STATUS.DISABLED ? 'danger' : 'success';
}

function handleUnavailableAction(action: 'create' | 'edit' | 'more') {
  MessagePlugin.warning(t(`user.userList.unavailable.${action}`));
}

function closeUserRoleDrawer() {
  drawerSession.value += 1;
  userRoleDrawerVisible.value = false;
  selectedUser.value = null;
  selectedRoleIds.value = [];
  loadingRoleSelection.value = false;
  roleSelectionReady.value = false;
  roleLoadWarning.value = '';
  submittingRoles.value = false;
}

function isActiveDrawerSession(session: number) {
  return userRoleDrawerVisible.value && drawerSession.value === session;
}

async function loadUserRoleSelection(user: UserRow, session: number) {
  selectedUser.value = user;
  selectedRoleIds.value = [];
  roleSelectionReady.value = false;
  roleLoadWarning.value = '';

  try {
    await loadRoleCatalog();
  } catch (error) {
    if (isActiveDrawerSession(session)) {
      roleLoadWarning.value = error instanceof Error ? error.message : t('user.userList.roleDialog.roleLoadFailed');
    }
    return;
  }

  if (!isActiveDrawerSession(session)) {
    return;
  }

  loadingRoleSelection.value = true;
  try {
    const response = await getUserRoleBindings(user.id);

    if (!isActiveDrawerSession(session)) {
      return;
    }

    selectedRoleIds.value = response.role_ids;
    roleSelectionReady.value = true;
  } catch (error) {
    if (isActiveDrawerSession(session)) {
      roleLoadWarning.value =
        error instanceof Error ? error.message : t('user.userList.roleDialog.selectionLoadFailed');
    }
  } finally {
    if (isActiveDrawerSession(session)) {
      loadingRoleSelection.value = false;
    }
  }
}

async function handleOpenUserRoleDrawer(row: UserRow) {
  const session = drawerSession.value + 1;

  drawerSession.value = session;
  userRoleDrawerVisible.value = true;
  await loadUserRoleSelection(row, session);
}

async function retryUserRoleDrawerLoad() {
  if (!selectedUser.value) {
    return;
  }

  await loadUserRoleSelection(selectedUser.value, drawerSession.value);
}

async function submitUserRoleAssignment() {
  if (!selectedUser.value || !canSubmitRoleAssignment.value) {
    return;
  }

  const session = drawerSession.value;
  submittingRoles.value = true;

  try {
    await assignUserRoles(selectedUser.value.id, {
      role_ids: [...selectedRoleIds.value].sort((left, right) => left - right),
    });

    if (!isActiveDrawerSession(session)) {
      return;
    }

    roleBindings.value = {
      ...roleBindings.value,
      [selectedUser.value.id]: [...selectedRoleIds.value],
    };
    MessagePlugin.success(t('user.userList.assignSuccess'));
    closeUserRoleDrawer();
  } catch (error) {
    if (isActiveDrawerSession(session)) {
      MessagePlugin.error(error instanceof Error ? error.message : t('user.userList.assignFailed'));
    }
  } finally {
    if (drawerSession.value === session) {
      submittingRoles.value = false;
    }
  }
}

onMounted(() => {
  fetchUsers();
});
</script>
<style lang="less" scoped>
@import './index.less';
</style>
