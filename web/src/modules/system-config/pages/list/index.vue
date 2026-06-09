<!--
  Copyright (c) 2025-2026 GeWuYou
  SPDX-License-Identifier: Apache-2.0
-->

<template>
  <section class="system-config-page" data-page-type="settings">
    <page-header
      :breadcrumb="[
        { labelKey: 'menu.server.title', fallback: t('systemConfig.list.eyebrow') },
        { labelKey: 'systemConfig.list.title', fallback: t('systemConfig.list.title') },
      ]"
      :source="{
        labelKey: 'menu.server.title',
        fallback: t('systemConfig.list.eyebrow'),
        color: 'var(--td-brand-color-6)',
      }"
      title-key="systemConfig.list.title"
      :title-fallback="t('systemConfig.list.title')"
      description-key="systemConfig.list.description"
      :description-fallback="t('systemConfig.list.description')"
    >
      <template #actions>
        <t-button theme="primary" :loading="loading" @click="refreshConfigs">
          <template #icon><refresh-icon /></template>
          {{ t('systemConfig.list.refresh') }}
        </t-button>
      </template>
    </page-header>

    <t-alert
      v-if="errorMessage"
      theme="error"
      :title="t('systemConfig.list.loadError')"
      :message="errorMessage"
      class="system-config-page__alert"
    />

    <t-loading :loading="loading">
      <div class="system-config-layout">
        <aside class="system-config-groups">
          <t-tree
            :data="domainTree"
            :actived="activeTreeValue"
            :expanded="expandedDomainKeys"
            activable
            hover
            expand-on-click-node
            @active="handleTreeActive"
          >
            <template #label="{ node }">
              <span class="system-config-tree-node">
                <span>{{ node.data.label }}</span>
                <small v-if="node.data.count">
                  {{ t('systemConfig.list.groupConfigCount', { count: node.data.count }) }}
                </small>
              </span>
            </template>
          </t-tree>
        </aside>

        <main class="system-config-content">
          <div v-if="activeGroup" class="system-config-content__head">
            <div>
              <h2>{{ activeGroup.label }}</h2>
              <p>{{ activeGroup.description }}</p>
            </div>
            <t-space size="small" break-line>
              <t-tag variant="light">
                {{ t('systemConfig.list.groupConfigCount', { count: activeGroup.items.length }) }}
              </t-tag>
              <t-tag :theme="activeGroupOverrideCount > 0 ? 'primary' : 'default'" variant="light">
                {{ t('systemConfig.list.overrideCount', { count: activeGroupOverrideCount }) }}
              </t-tag>
            </t-space>
          </div>

          <div v-if="activeGroup?.items.length" class="system-config-list">
            <t-card v-for="item in activeGroup.items" :key="item.key" class="system-config-item" bordered>
              <div class="system-config-item__main">
                <div class="system-config-item__title-row">
                  <div>
                    <h3>{{ configTitle(item) }}</h3>
                    <p>{{ configDescription(item) }}</p>
                  </div>
                  <t-space size="small" break-line>
                    <t-tag :theme="configStatus(item).theme" variant="light">
                      {{ configStatus(item).label }}
                    </t-tag>
                    <t-tag v-if="item.sensitive" theme="danger" variant="light">
                      {{ t('systemConfig.list.tags.sensitive') }}
                    </t-tag>
                    <t-tag v-if="item.restart_required" theme="primary" variant="light">
                      {{ t('systemConfig.list.tags.restartRequired') }}
                    </t-tag>
                  </t-space>
                </div>

                <div class="system-config-values">
                  <section
                    v-for="valueSection in valueSections(item)"
                    :key="valueSection.field"
                    class="system-config-value"
                  >
                    <header>
                      <h4>{{ valueSection.title }}</h4>
                    </header>
                    <dl class="system-config-value__rows">
                      <template v-for="row in valueSection.rows" :key="row.key">
                        <dt>{{ row.label }}</dt>
                        <dd>
                          <strong>{{ row.value }}</strong>
                          <small v-if="row.description">{{ row.description }}</small>
                        </dd>
                      </template>
                    </dl>
                    <t-collapse
                      v-if="valueSection.json"
                      borderless
                      expand-icon-placement="right"
                      class="system-config-json"
                    >
                      <t-collapse-panel :value="valueSection.field" :header="t('systemConfig.list.viewJson')">
                        <pre>{{ valueSection.json }}</pre>
                      </t-collapse-panel>
                    </t-collapse>
                  </section>
                </div>

                <div class="system-config-summary">
                  <section class="system-config-summary__cell">
                    <span>{{ t('systemConfig.list.lastModified.title') }}</span>
                    <strong>{{ configLastModifiedLabel(item) }}</strong>
                  </section>
                </div>

                <section class="system-config-technical">
                  <span>{{ t('systemConfig.list.technicalId') }}</span>
                  <code>{{ item.key }}</code>
                </section>

                <div class="system-config-item__actions">
                  <t-button
                    v-permission="permissionCodes.WRITE"
                    theme="primary"
                    variant="outline"
                    @click="openEditor(item)"
                  >
                    <template #icon><edit-icon /></template>
                    {{ t('systemConfig.list.edit') }}
                  </t-button>
                  <t-popconfirm
                    v-if="item.has_override"
                    theme="warning"
                    :content="t('systemConfig.list.resetConfirm')"
                    :confirm-btn="t('systemConfig.list.reset')"
                    :cancel-btn="t('systemConfig.list.cancel')"
                    @confirm="resetConfigOverride(item)"
                  >
                    <t-button
                      v-permission="permissionCodes.WRITE"
                      theme="default"
                      variant="outline"
                      :loading="resettingKey === item.key"
                    >
                      <template #icon><rollback-icon /></template>
                      {{ t('systemConfig.list.reset') }}
                    </t-button>
                  </t-popconfirm>
                </div>
              </div>
            </t-card>
          </div>

          <t-empty
            v-else
            :title="t('systemConfig.list.emptyTitle')"
            :description="t('systemConfig.list.emptyDescription')"
          >
            <template #action>
              <t-button theme="primary" variant="outline" @click="refreshConfigs">
                {{ t('systemConfig.list.refresh') }}
              </t-button>
            </template>
          </t-empty>
        </main>
      </div>
    </t-loading>

    <t-dialog
      v-model:visible="editorVisible"
      :header="editorTitle"
      :confirm-btn="t('systemConfig.list.save')"
      :cancel-btn="t('systemConfig.list.cancel')"
      :confirm-loading="saving"
      width="680px"
      destroy-on-close
      @confirm="saveEditor"
    >
      <div v-if="editingItem" class="system-config-editor">
        <t-alert v-if="editingItem.sensitive" theme="warning" :message="t('systemConfig.list.sensitiveEditHint')" />
        <t-form :data="editorForm" label-align="top">
          <json-schema-value-fields
            v-model="editorForm.value"
            :root-schema="editingSchema"
            :labels="schemaLabels"
            :title-resolver="schemaFieldTitle"
            :description-resolver="schemaFieldDescription"
            :placeholder-resolver="schemaFieldPlaceholder"
            :unit-resolver="schemaFieldUnit"
            :option-label-resolver="schemaOptionLabel"
          />
        </t-form>
        <section class="system-config-editor__preview">
          <strong>{{ t('systemConfig.list.previewTitle') }}</strong>
          <pre>{{ editorPreview }}</pre>
        </section>
      </div>
    </t-dialog>
  </section>
</template>
<script setup lang="ts">
import { EditIcon, RefreshIcon, RollbackIcon } from 'tdesign-icons-vue-next';
import type { TreeNodeValue, TreeProps } from 'tdesign-vue-next';
import { MessagePlugin } from 'tdesign-vue-next';
import { computed, onMounted, reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { formatCompactDateTime } from '@/shared/components/management';
import { PageHeader } from '@/shared/components/page';
import {
  type ConfigSchemaField,
  getConfigSchemaFields,
  JsonSchemaValueFields,
  parseConfigSchema,
} from '@/shared/schema-form';
import { formatJsonValue, isJsonRecord, parseJsonValue, valuePreview } from '@/shared/schema-form/json';
import type { ApiRequestError } from '@/types/axios';

import { getSystemConfigs, resetSystemConfig, updateSystemConfig } from '../../api/system-config';
import { SYSTEM_CONFIG_PERMISSION_CODE } from '../../contract/permissions';
import type { SystemConfigItem } from '../../types/system-config';

defineOptions({
  name: 'SystemConfigListPage',
});

type ConfigGroup = {
  key: string;
  domainKey: string;
  label: string;
  description: string;
  items: SystemConfigItem[];
};

type ConfigDomain = {
  key: string;
  label: string;
  groups: ConfigGroup[];
};

type ConfigValueField = 'effective_value' | 'default_value';

type ConfigValueRow = {
  key: string;
  label: string;
  description: string;
  value: string;
};

type ConfigValueSection = {
  field: ConfigValueField;
  title: string;
  rows: ConfigValueRow[];
  json: string;
};

const { locale, t, te } = useI18n();
const permissionCodes = SYSTEM_CONFIG_PERMISSION_CODE;
const items = ref<SystemConfigItem[]>([]);
const loading = ref(false);
const saving = ref(false);
const resettingKey = ref('');
const errorMessage = ref('');
const activeGroupKey = ref('');
const expandedDomainKeys = ref<TreeNodeValue[]>([]);
const editorVisible = ref(false);
const editingItem = ref<SystemConfigItem | null>(null);
const editorForm = reactive<{ value: unknown }>({ value: undefined });

const schemaLabels = computed(() => ({
  invalidJson: t('systemConfig.list.schema.invalidJson'),
  jsonPlaceholder: t('systemConfig.list.schema.jsonPlaceholder'),
  numberPlaceholder: t('systemConfig.list.schema.numberPlaceholder'),
  selectPlaceholder: t('systemConfig.list.schema.selectPlaceholder'),
  stringPlaceholder: t('systemConfig.list.schema.stringPlaceholder'),
  value: t('systemConfig.list.schema.value'),
}));

const domains = computed<ConfigDomain[]>(() => {
  const domainMap = new Map<string, ConfigDomain>();
  const sortedItems = [...items.value].sort((left, right) => {
    const orderDelta = (left.order ?? 0) - (right.order ?? 0);
    return orderDelta || configTitle(left).localeCompare(configTitle(right));
  });

  for (const item of sortedItems) {
    const domainKey = configDomainKey(item);
    const domain = domainMap.get(domainKey) ?? {
      key: domainKey,
      label: domainLabel(item),
      groups: [],
    };
    const groupKey = `${domainKey}:${item.module}:${item.group || 'default'}`;
    let group = domain.groups.find((candidate) => candidate.key === groupKey);
    if (!group) {
      group = {
        key: groupKey,
        domainKey,
        label: groupLabel(item),
        description: groupDescription(item),
        items: [],
      };
      domain.groups.push(group);
    }
    group.items.push(item);
    domainMap.set(domainKey, domain);
  }

  return [...domainMap.values()];
});

const groupedConfigs = computed(() => domains.value.flatMap((domain) => domain.groups));
const domainTree = computed<TreeProps['data']>(() =>
  domains.value.map((domain) => ({
    value: domain.key,
    label: domain.label,
    children: domain.groups.map((group) => ({
      value: group.key,
      label: group.label,
      count: group.items.length,
    })),
  })),
);
const activeTreeValue = computed(() => (activeGroupKey.value ? [activeGroupKey.value] : []));
const activeGroup = computed(() => groupedConfigs.value.find((group) => group.key === activeGroupKey.value) ?? null);
const activeGroupOverrideCount = computed(
  () => activeGroup.value?.items.filter((item) => item.has_override).length ?? 0,
);
const editingSchema = computed(() => parseConfigSchema(editingItem.value?.config_schema));
const editorTitle = computed(() =>
  editingItem.value ? t('systemConfig.list.editorTitle', { title: configTitle(editingItem.value) }) : '',
);
const editorPreview = computed(() => formatJsonValue(editorForm.value) || t('systemConfig.list.emptyValue'));

onMounted(refreshConfigs);

async function refreshConfigs() {
  loading.value = true;
  errorMessage.value = '';
  try {
    const response = await getSystemConfigs();
    items.value = response.items ?? [];
    if (!activeGroupKey.value || !groupedConfigs.value.some((group) => group.key === activeGroupKey.value)) {
      activeGroupKey.value = groupedConfigs.value[0]?.key ?? '';
    }
    expandedDomainKeys.value = domains.value.map((domain) => domain.key);
  } catch (error) {
    errorMessage.value = readableError(error, t('systemConfig.list.loadError'));
  } finally {
    loading.value = false;
  }
}

function handleTreeActive(value: TreeNodeValue[]) {
  const selected = String(value[0] ?? '');
  if (groupedConfigs.value.some((group) => group.key === selected)) {
    activeGroupKey.value = selected;
  }
}

function openEditor(item: SystemConfigItem) {
  editingItem.value = item;
  editorForm.value = initialEditorValue(item);
  editorVisible.value = true;
}

async function saveEditor() {
  if (!editingItem.value) {
    return;
  }

  saving.value = true;
  try {
    const updated = await updateSystemConfig(editingItem.value.key, { value: editorForm.value });
    upsertConfig(updated);
    editingItem.value = updated;
    editorVisible.value = false;
    MessagePlugin.success(t('systemConfig.list.saveSuccess'));
  } catch (error) {
    MessagePlugin.error(readableError(error, t('systemConfig.list.saveError')));
  } finally {
    saving.value = false;
  }
}

async function resetConfigOverride(item: SystemConfigItem) {
  resettingKey.value = item.key;
  try {
    const updated = await resetSystemConfig(item.key);
    upsertConfig(updated);
    MessagePlugin.success(t('systemConfig.list.resetSuccess'));
  } catch (error) {
    MessagePlugin.error(readableError(error, t('systemConfig.list.resetError')));
  } finally {
    resettingKey.value = '';
  }
}

function upsertConfig(item: SystemConfigItem) {
  const nextItems = [...items.value];
  const index = nextItems.findIndex((candidate) => candidate.key === item.key);
  if (index >= 0) {
    nextItems[index] = item;
  } else {
    nextItems.push(item);
  }
  items.value = nextItems;
}

function initialEditorValue(item: SystemConfigItem) {
  if (item.sensitive) {
    return emptySensitiveEditorValue(item.type);
  }

  return (
    parseJsonValue(item.override_value) ?? parseJsonValue(item.effective_value) ?? parseJsonValue(item.default_value)
  );
}

function emptySensitiveEditorValue(type: SystemConfigItem['type']) {
  switch (type) {
    case 'boolean':
      return false;
    case 'object':
      return {};
    case 'array':
      return [];
    case 'number':
    case 'integer':
      return null;
    default:
      return '';
  }
}

function configTitle(item: SystemConfigItem) {
  return resolveI18nText(item.title_key, item.title, item.key);
}

function configDescription(item: SystemConfigItem) {
  return resolveI18nText(item.description_key, item.description, t('systemConfig.list.noDescription'));
}

function groupLabel(item: SystemConfigItem) {
  return resolveI18nText(item.group_key, item.group_label, technicalGroupKey(item));
}

function groupDescription(item: SystemConfigItem) {
  return resolveI18nText(item.group_description_key, item.group_description, t('systemConfig.list.noDescription'));
}

function domainLabel(item: SystemConfigItem) {
  return resolveI18nText(item.domain_key, item.domain_label, technicalDomainKey(item));
}

function configDomainKey(item: SystemConfigItem) {
  return item.domain?.trim() || item.module || 'uncategorized';
}

function technicalDomainKey(item: SystemConfigItem) {
  return item.domain?.trim() || t('systemConfig.domains.uncategorized');
}

function technicalGroupKey(item: SystemConfigItem) {
  return item.group || t('systemConfig.list.defaultGroup');
}

function configStatus(item: SystemConfigItem) {
  if (isModifiedConfig(item)) {
    return {
      label: t('systemConfig.list.status.modified'),
      description: t('systemConfig.list.status.modifiedDescription'),
      theme: 'primary' as const,
    };
  }

  return {
    label: t('systemConfig.list.status.default'),
    description: t('systemConfig.list.status.defaultDescription'),
    theme: 'default' as const,
  };
}

function isModifiedConfig(item: SystemConfigItem) {
  return item.status === 'modified';
}

function configLastModifiedLabel(item: SystemConfigItem) {
  if (!isModifiedConfig(item)) {
    return t('systemConfig.list.lastModified.none');
  }

  const updatedAt = formatCompactDateTime(item.updated_at, locale);
  const userLabel = configUpdatedByLabel(item);
  return t('systemConfig.list.lastModified.value', { user: userLabel, time: updatedAt });
}

function configUpdatedByLabel(item: SystemConfigItem) {
  const username = item.updated_by_username?.trim();
  if (username) {
    return username;
  }

  if (item.updated_by_user_id !== undefined && item.updated_by_user_id !== null) {
    return t('systemConfig.list.lastModified.userId', { id: item.updated_by_user_id });
  }

  return t('systemConfig.list.lastModified.unknownUser');
}

function valueSections(item: SystemConfigItem): ConfigValueSection[] {
  const sections = [buildValueSection(item, 'effective_value', t('systemConfig.list.values.current'))];
  if (item.has_override) {
    sections.push(buildValueSection(item, 'default_value', t('systemConfig.list.values.default')));
  }
  return sections;
}

function buildValueSection(item: SystemConfigItem, field: ConfigValueField, title: string): ConfigValueSection {
  if (item.masked) {
    const maskedValue = item.masked_placeholder || t('systemConfig.list.masked');
    return {
      field,
      title,
      rows: [{ key: field, label: title, description: '', value: maskedValue }],
      json: '',
    };
  }

  const parsed = parseJsonValue(item[field]);
  const schema = parseConfigSchema(item.config_schema);
  const fields = getConfigSchemaFields(schema);
  const rows = isJsonRecord(parsed) && fields.length > 0 ? structuredValueRows(parsed, fields) : [];
  const fallbackRows =
    rows.length > 0
      ? rows
      : [
          {
            key: field,
            label: schema.title
              ? resolveI18nText(schema.xI18n?.titleKey, schema.title, configTitle(item))
              : configTitle(item),
            description: resolveI18nText(schema.xI18n?.descriptionKey, schema.description, ''),
            value: configValuePreview(parsed, schema),
          },
        ];

  return {
    field,
    title,
    rows: fallbackRows,
    json: shouldShowJsonValue(parsed) ? formatJsonValue(parsed) : '',
  };
}

function structuredValueRows(value: Record<string, unknown>, fields: ConfigSchemaField[]): ConfigValueRow[] {
  return fields.map((field) => {
    const unit = schemaFieldUnit(field);
    const displayValue = valuePreview(value[field.key], t('systemConfig.list.emptyValue'), booleanStateLabel);
    return {
      key: field.key,
      label: schemaFieldTitle(field),
      description: schemaFieldDescription(field),
      value: unit ? `${displayValue} ${unit}` : displayValue,
    };
  });
}

function booleanStateLabel(value: boolean) {
  return value ? t('systemConfig.list.boolean.enabled') : t('systemConfig.list.boolean.disabled');
}

function shouldShowJsonValue(value: unknown) {
  return Array.isArray(value) || isJsonRecord(value);
}

function configValuePreview(value: unknown, schema = parseConfigSchema()) {
  const optionText = schema.enumLabels?.[String(value)];
  if (optionText) {
    return resolveI18nText(optionText.labelKey, optionText.label, String(value));
  }
  return valuePreview(value, t('systemConfig.list.emptyValue'), booleanStateLabel);
}

function schemaFieldTitle(field: ConfigSchemaField) {
  return resolveI18nText(field.schema.xI18n?.titleKey, field.schema.title, field.key);
}

function schemaFieldDescription(field: ConfigSchemaField) {
  return resolveI18nText(field.schema.xI18n?.descriptionKey, field.schema.description, '');
}

function schemaFieldPlaceholder(field: ConfigSchemaField) {
  return resolveI18nText(field.schema.xI18n?.placeholderKey, field.schema.placeholder, '');
}

function schemaFieldUnit(field: ConfigSchemaField) {
  return resolveI18nText(field.schema.xI18n?.unitKey, undefined, '');
}

function schemaOptionLabel(field: ConfigSchemaField, option: string | number | boolean) {
  const optionText = field.schema.enumLabels?.[String(option)];
  return resolveI18nText(optionText?.labelKey, optionText?.label, String(option));
}

function resolveI18nText(key?: string, fallback?: string, rawFallback = '') {
  return key && te(key) ? t(key) : fallback || rawFallback;
}

function readableError(error: unknown, fallback: string) {
  const requestError = error as Partial<ApiRequestError>;
  return typeof requestError.message === 'string' && requestError.message.trim() ? requestError.message : fallback;
}
</script>
<style scoped>
.system-config-page {
  color: var(--td-text-color-primary);
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-16);
  min-width: 0;
}

.system-config-content__head h2,
.system-config-item h3 {
  margin: 0;
}

.system-config-content__head p,
.system-config-item p {
  color: var(--td-text-color-secondary);
  margin: var(--graft-density-gap-4) 0 0;
}

.system-config-layout {
  align-items: flex-start;
  display: grid;
  gap: var(--graft-density-gap-16);
  grid-template-columns: minmax(220px, 280px) minmax(0, 1fr);
}

.system-config-groups {
  background: var(--td-bg-color-container);
  border: 1px solid var(--td-border-level-1-color);
  border-radius: var(--td-radius-medium);
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-8);
  padding: var(--graft-density-gap-12);
  position: sticky;
  top: var(--graft-density-gap-16);
}

.system-config-groups :deep(.t-tree) {
  background: transparent;
}

.system-config-tree-node {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-4);
  line-height: 1.4;
  min-width: 0;
}

.system-config-tree-node small {
  color: var(--td-text-color-secondary);
}

.system-config-content {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-12);
  min-width: 0;
}

.system-config-content__head,
.system-config-item__title-row {
  align-items: flex-start;
  display: flex;
  gap: var(--graft-density-gap-12);
  justify-content: space-between;
}

.system-config-list {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-12);
}

.system-config-item {
  background: var(--td-bg-color-container);
  border: 1px solid var(--td-border-level-1-color);
  border-radius: var(--td-radius-medium);
}

.system-config-item :deep(.t-card__body) {
  padding: var(--graft-density-gap-16);
}

.system-config-item__main,
.system-config-editor {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-12);
  min-width: 0;
}

.system-config-item__actions {
  display: flex;
  flex-flow: row wrap;
  gap: var(--graft-density-gap-8);
  justify-content: flex-end;
}

.system-config-summary,
.system-config-values {
  display: grid;
  gap: var(--graft-density-gap-8);
}

.system-config-summary {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.system-config-summary__cell,
.system-config-value,
.system-config-technical {
  background: var(--td-bg-color-page);
  border-radius: var(--td-radius-small);
  padding: var(--graft-density-gap-12);
}

.system-config-summary__cell {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-8);
}

.system-config-summary__cell > span,
.system-config-value h4,
.system-config-technical span {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-small);
  margin: 0;
}

.system-config-summary__cell > div {
  align-items: center;
  display: flex;
  flex-wrap: wrap;
  gap: var(--graft-density-gap-8);
}

.system-config-values {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.system-config-value {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-12);
  min-width: 0;
}

.system-config-value__rows {
  display: grid;
  gap: var(--graft-density-gap-8) var(--graft-density-gap-16);
  grid-template-columns: minmax(120px, max-content) minmax(0, 1fr);
  margin: 0;
}

.system-config-value__rows dt {
  color: var(--td-text-color-secondary);
}

.system-config-value__rows dd {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-4);
  margin: 0;
  min-width: 0;
}

.system-config-value__rows strong,
.system-config-technical code {
  overflow-wrap: anywhere;
}

.system-config-value__rows small {
  color: var(--td-text-color-placeholder);
}

.system-config-json :deep(.t-collapse-panel__content) {
  padding: 0;
}

.system-config-technical {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-8);
}

.system-config-json pre,
.system-config-editor__preview pre {
  border-radius: var(--td-radius-small);
  box-sizing: border-box;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', monospace;
  margin: var(--graft-density-gap-8) 0 0;
  max-height: 220px;
  overflow: auto;
  padding: var(--graft-density-gap-12);
  white-space: pre-wrap;
}

.system-config-json pre {
  background: var(--td-bg-color-container);
  border: 1px solid var(--td-border-level-1-color);
}

.system-config-editor__preview pre {
  background: var(--td-bg-color-page);
}

@media (width <= 900px) {
  .system-config-layout,
  .system-config-summary,
  .system-config-values {
    display: flex;
    flex-direction: column;
  }

  .system-config-groups {
    position: static;
    width: 100%;
  }

  .system-config-item__actions {
    flex-flow: row wrap;
  }
}
</style>
