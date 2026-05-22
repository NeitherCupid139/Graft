<template>
  <div class="access-control-overview" data-page-type="overview-dashboard">
    <management-page-header
      :title="t('accessControl.overview.title')"
      :description="t('accessControl.overview.description')"
    >
      <template #eyebrow>{{ t('accessControl.overview.navHint') }}</template>
      <template #actions>
        <t-button theme="primary" @click="goToUsers">{{ t('accessControl.overview.actions.createUser') }}</t-button>
        <t-button variant="outline" @click="goToRoles">{{ t('accessControl.overview.actions.createRole') }}</t-button>
        <t-button variant="outline" @click="goToPermissions">{{
          t('accessControl.overview.actions.viewPermissions')
        }}</t-button>
        <t-button variant="outline" :loading="loading" @click="fetchOverview">
          {{ t('accessControl.overview.actions.refresh') }}
        </t-button>
      </template>
    </management-page-header>

    <div v-if="loadError" class="access-control-overview__feedback">
      <management-empty-state
        tone="error"
        :title="t('accessControl.overview.state.loadFailedTitle')"
        :description="loadError || t('accessControl.overview.state.loadFailedDescription')"
      >
        <template #actions>
          <t-button theme="primary" variant="outline" @click="fetchOverview">
            {{ t('accessControl.overview.state.retry') }}
          </t-button>
        </template>
      </management-empty-state>
    </div>

    <management-stats-grid :items="statItems" />

    <section class="access-control-overview__grid">
      <management-table-card>
        <template #head>
          <div class="section-head">
            <h2>{{ t('accessControl.overview.risk.title') }}</h2>
          </div>
        </template>
        <div v-if="riskItems.length > 0" class="risk-list">
          <article v-for="item in riskItems" :key="item.label" class="risk-list__item">
            <span class="risk-list__label">{{ item.label }}</span>
            <strong class="risk-list__value">{{ item.value }}</strong>
            <p v-if="item.description" class="risk-list__description">{{ item.description }}</p>
          </article>
        </div>
        <management-empty-state
          v-else
          :title="t('accessControl.overview.risk.emptyTitle')"
          :description="t('accessControl.overview.risk.emptyDescription')"
        />
      </management-table-card>

      <management-table-card>
        <template #head>
          <div class="section-head">
            <h2>{{ t('accessControl.overview.quickLinks.title') }}</h2>
          </div>
        </template>
        <div class="quick-link-grid">
          <button class="quick-link-card" type="button" @click="goToUsers">
            <span class="quick-link-card__title">{{ t('accessControl.overview.quickLinks.users.title') }}</span>
            <span class="quick-link-card__count">{{ displayValue(users.length) }}</span>
            <p class="quick-link-card__description">{{ t('accessControl.overview.quickLinks.users.description') }}</p>
          </button>
          <button class="quick-link-card" type="button" @click="goToRoles">
            <span class="quick-link-card__title">{{ t('accessControl.overview.quickLinks.roles.title') }}</span>
            <span class="quick-link-card__count">{{ displayValue(roles.length) }}</span>
            <p class="quick-link-card__description">{{ t('accessControl.overview.quickLinks.roles.description') }}</p>
          </button>
          <button class="quick-link-card" type="button" @click="goToPermissions">
            <span class="quick-link-card__title">{{ t('accessControl.overview.quickLinks.permissions.title') }}</span>
            <span class="quick-link-card__count">{{ displayValue(permissions.length) }}</span>
            <p class="quick-link-card__description">
              {{ t('accessControl.overview.quickLinks.permissions.description') }}
            </p>
          </button>
        </div>
      </management-table-card>
    </section>

    <management-table-card>
      <template #head>
        <div class="section-head">
          <h2>{{ t('accessControl.overview.changes.title') }}</h2>
        </div>
      </template>
      <management-empty-state
        :title="t('accessControl.overview.changes.emptyTitle')"
        :description="t('accessControl.overview.changes.emptyDescription')"
      />
    </management-table-card>
  </div>
</template>
<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import { getPermissions } from '@/modules/rbac/api/rbac';
import { getRoles as getUserRoles } from '@/modules/user/api/user-roles';
import { getUsers } from '@/modules/user/api/users';
import {
  ManagementEmptyState,
  ManagementPageHeader,
  type ManagementStatItem,
  ManagementStatsGrid,
  ManagementTableCard,
} from '@/shared/components/management';

import { ACCESS_CONTROL_ROUTE_PATH } from '../../contract/bootstrap';

type RoleSummary = {
  id: number;
  builtin: boolean;
  permission_count: number;
  updated_at: string;
};

type PermissionSummary = {
  id: number;
};

type UserSummary = {
  id: number;
  status: string;
};

const { t } = useI18n();
const router = useRouter();
const loading = ref(false);
const loadError = ref('');
const users = ref<UserSummary[]>([]);
const roles = ref<RoleSummary[]>([]);
const permissions = ref<PermissionSummary[]>([]);
const roleBindings = ref<Record<number, number[]>>({});

const unassignedUserCount = computed(
  () => users.value.filter((user) => (roleBindings.value[user.id] ?? []).length === 0).length,
);
const disabledUserCount = computed(() => users.value.filter((user) => user.status === 'disabled').length);
const builtinRoleCount = computed(() => roles.value.filter((role) => role.builtin).length);
const customRoleCount = computed(() => roles.value.filter((role) => !role.builtin).length);
const emptyCustomRoleCount = computed(
  () => roles.value.filter((role) => !role.builtin && role.permission_count === 0).length,
);
const assignedUserCount = computed(
  () => users.value.filter((user) => (roleBindings.value[user.id] ?? []).length > 0).length,
);
const staleRoleCount = computed(() => roles.value.filter((role) => !role.updated_at).length);

const statItems = computed<ManagementStatItem[]>(() => [
  { label: t('accessControl.overview.stats.totalUsers'), value: displayValue(users.value.length) },
  {
    label: t('accessControl.overview.stats.enabledUsers'),
    value: displayValue(users.value.length - disabledUserCount.value),
  },
  { label: t('accessControl.overview.stats.disabledUsers'), value: displayValue(disabledUserCount.value) },
  { label: t('accessControl.overview.stats.totalRoles'), value: displayValue(roles.value.length) },
  { label: t('accessControl.overview.stats.builtinRoles'), value: displayValue(builtinRoleCount.value) },
  { label: t('accessControl.overview.stats.customRoles'), value: displayValue(customRoleCount.value) },
  { label: t('accessControl.overview.stats.totalPermissions'), value: displayValue(permissions.value.length) },
  { label: t('accessControl.overview.stats.assignedUsers'), value: displayValue(assignedUserCount.value) },
]);

const riskItems = computed(() => {
  const items: Array<{ label: string; value: string; description?: string }> = [];

  if (unassignedUserCount.value > 0) {
    items.push({
      label: t('accessControl.overview.risk.unassignedUsers'),
      value: displayValue(unassignedUserCount.value),
    });
  }

  if (emptyCustomRoleCount.value > 0) {
    items.push({
      label: t('accessControl.overview.risk.emptyRoles'),
      value: displayValue(emptyCustomRoleCount.value),
    });
  }

  if (disabledUserCount.value > 0) {
    items.push({
      label: t('accessControl.overview.risk.disabledUsers'),
      value: displayValue(disabledUserCount.value),
    });
  }

  if (roles.value.some((role) => role.builtin)) {
    items.push({
      label: t('accessControl.overview.risk.builtinNotice'),
      value: displayValue(builtinRoleCount.value),
      description: t('accessControl.overview.risk.builtinNoticeDescription'),
    });
  }

  if (staleRoleCount.value > 0) {
    items.push({
      label: t('accessControl.overview.risk.stale'),
      value: displayValue(staleRoleCount.value),
    });
  }

  return items;
});

function displayValue(value?: number | null) {
  return typeof value === 'number' && Number.isFinite(value)
    ? new Intl.NumberFormat().format(value)
    : t('accessControl.overview.state.unknown');
}

async function fetchOverview() {
  loading.value = true;
  loadError.value = '';

  try {
    const [userResult, roleResult, permissionResult] = await Promise.all([
      getUsers(),
      getUserRoles(),
      getPermissions(),
    ]);
    users.value = userResult.items;
    roles.value = roleResult.items;
    permissions.value = permissionResult.items;

    const bindings = await Promise.all(
      userResult.items.map(async (user) => {
        const { getUserRoleBindings } = await import('@/modules/user/api/user-roles');
        const response = await getUserRoleBindings(user.id);
        return [user.id, response.role_ids] as const;
      }),
    );

    roleBindings.value = Object.fromEntries(bindings);
  } catch (error) {
    users.value = [];
    roles.value = [];
    permissions.value = [];
    roleBindings.value = {};
    loadError.value = error instanceof Error ? error.message : t('accessControl.overview.state.loadFailedDescription');
  } finally {
    loading.value = false;
  }
}

function goToUsers() {
  void router.push(ACCESS_CONTROL_ROUTE_PATH.USERS);
}

function goToRoles() {
  void router.push(ACCESS_CONTROL_ROUTE_PATH.ROLES);
}

function goToPermissions() {
  void router.push(ACCESS_CONTROL_ROUTE_PATH.PERMISSIONS);
}

onMounted(() => {
  fetchOverview();
});
</script>
<style scoped lang="less">
.access-control-overview {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.access-control-overview__grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.section-head h2 {
  color: var(--td-text-color-primary);
  font: var(--td-font-title-large);
  margin: 0;
}

.risk-list {
  display: grid;
  gap: 12px;
}

.risk-list__item {
  background: var(--td-bg-color-secondarycontainer);
  border: 1px solid var(--td-component-border);
  border-radius: var(--td-radius-medium);
  display: grid;
  gap: 6px;
  padding: 14px 16px;
}

.risk-list__label,
.risk-list__description,
.quick-link-card__description {
  color: var(--td-text-color-secondary);
}

.risk-list__value,
.quick-link-card__count {
  color: var(--td-text-color-primary);
  font: var(--td-font-title-large);
}

.risk-list__description {
  font: var(--td-font-body-small);
  margin: 0;
}

.quick-link-grid {
  display: grid;
  gap: 12px;
}

.quick-link-card {
  background: var(--td-bg-color-secondarycontainer);
  border: 1px solid var(--td-component-border);
  border-radius: var(--td-radius-large);
  color: inherit;
  cursor: pointer;
  display: grid;
  gap: 8px;
  padding: 16px;
  text-align: left;
  transition:
    border-color 0.2s ease,
    transform 0.2s ease;
}

.quick-link-card:hover {
  border-color: var(--td-brand-color);
  transform: translateY(-1px);
}

.quick-link-card__title {
  color: var(--td-text-color-primary);
  font: var(--td-font-title-medium);
}

.quick-link-card__description {
  margin: 0;
}

@media (width <= 900px) {
  .access-control-overview__grid {
    grid-template-columns: 1fr;
  }
}
</style>
