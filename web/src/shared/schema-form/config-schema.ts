import { isJsonRecord, type JsonRecord, parseJsonRecord } from './json';

export type ConfigFieldType = 'string' | 'integer' | 'number' | 'boolean' | 'object' | 'array';

export type ConfigSchemaProperty = {
  type?: ConfigFieldType;
  enum?: Array<string | number | boolean>;
  enumLabels?: Record<string, ConfigSchemaOptionLabel>;
  minimum?: number;
  maximum?: number;
  minLength?: number;
  maxLength?: number;
  default?: unknown;
  title?: string;
  description?: string;
  placeholder?: string;
  xI18n?: ConfigSchemaI18n;
  'x-title-key'?: string;
  'x-description-key'?: string;
};

export type ConfigSchemaI18n = {
  titleKey?: string;
  descriptionKey?: string;
  placeholderKey?: string;
  unitKey?: string;
};

export type ConfigSchemaOptionLabel = {
  label?: string;
  labelKey?: string;
};

export type ConfigSchema = ConfigSchemaProperty & {
  type?: ConfigFieldType;
  properties?: Record<string, ConfigSchemaProperty>;
  required?: string[];
  additionalProperties?: boolean;
};

export type ConfigSchemaField = {
  key: string;
  schema: ConfigSchemaProperty;
  required: boolean;
};

export function parseConfigSchema(value?: string | null | JsonRecord): ConfigSchema {
  const parsed = typeof value === 'string' || value === null || value === undefined ? parseJsonRecord(value) : value;
  return parseConfigSchemaRecord(parsed);
}

function parseConfigSchemaRecord(raw: JsonRecord): ConfigSchema {
  const type = parseFieldType(raw.type);

  return {
    ...parseProperty(raw),
    type,
    properties: parseProperties(raw.properties),
    required: Array.isArray(raw.required)
      ? raw.required.filter((item): item is string => typeof item === 'string')
      : [],
    additionalProperties: typeof raw.additionalProperties === 'boolean' ? raw.additionalProperties : undefined,
  };
}

export function getConfigSchemaFields(schema: ConfigSchema): ConfigSchemaField[] {
  const required = new Set(schema.required ?? []);
  return Object.entries(schema.properties ?? {}).map(([key, property]) => ({
    key,
    schema: property,
    required: required.has(key),
  }));
}

export function mergeConfigRecords(defaultConfig: JsonRecord, config: JsonRecord): JsonRecord {
  return {
    ...defaultConfig,
    ...config,
  };
}

export function buildDefaultConfigFromSchema(schema: ConfigSchema): JsonRecord {
  const output: JsonRecord = {};
  for (const [key, property] of Object.entries(schema.properties ?? {})) {
    if ('default' in property) {
      output[key] = property.default;
    }
  }
  return output;
}

function parseProperties(value: unknown): Record<string, ConfigSchemaProperty> {
  if (!isJsonRecord(value)) {
    return {};
  }

  return Object.fromEntries(
    Object.entries(value)
      .filter((entry): entry is [string, JsonRecord] => isJsonRecord(entry[1]))
      .map(([key, raw]) => [key, parseProperty(raw)]),
  );
}

function parseProperty(raw: JsonRecord): ConfigSchemaProperty {
  const type = parseFieldType(raw.type);
  const property: ConfigSchemaProperty = {};

  if (type) {
    property.type = type;
  }
  if (Array.isArray(raw.enum)) {
    property.enum = raw.enum.filter(
      (value): value is string | number | boolean =>
        typeof value === 'string' || typeof value === 'number' || typeof value === 'boolean',
    );
  }
  property.enumLabels = parseEnumLabels(raw);
  assignNumber(raw, property, 'minimum');
  assignNumber(raw, property, 'maximum');
  assignNumber(raw, property, 'minLength');
  assignNumber(raw, property, 'maxLength');
  assignString(raw, property, 'title');
  assignString(raw, property, 'description');
  assignString(raw, property, 'placeholder');
  assignString(raw, property, 'x-title-key');
  assignString(raw, property, 'x-description-key');
  property.xI18n = parseI18n(raw);
  if ('default' in raw) {
    property.default = raw.default;
  }

  return property;
}

function parseI18n(raw: JsonRecord): ConfigSchemaI18n | undefined {
  const extension = isJsonRecord(raw['x-i18n']) ? raw['x-i18n'] : raw;
  const parsed: ConfigSchemaI18n = {};
  assignI18nString(extension, parsed, 'titleKey');
  assignI18nString(extension, parsed, 'descriptionKey');
  assignI18nString(extension, parsed, 'placeholderKey');
  assignI18nString(extension, parsed, 'unitKey');
  assignLegacyI18nString(raw, parsed, 'x-title-key', 'titleKey');
  assignLegacyI18nString(raw, parsed, 'x-description-key', 'descriptionKey');
  return Object.keys(parsed).length > 0 ? parsed : undefined;
}

function parseEnumLabels(raw: JsonRecord): Record<string, ConfigSchemaOptionLabel> | undefined {
  const rawOptions = isJsonRecord(raw.options)
    ? raw.options
    : isJsonRecord(raw.enumLabels)
      ? raw.enumLabels
      : undefined;
  if (!rawOptions) {
    return undefined;
  }

  const entries = Object.entries(rawOptions).flatMap(([value, option]) => {
    if (typeof option === 'string') {
      return [[value, { label: option } satisfies ConfigSchemaOptionLabel]];
    }
    if (!isJsonRecord(option)) {
      return [];
    }
    const label: ConfigSchemaOptionLabel = {};
    assignOptionLabelString(option, label, 'label');
    assignOptionLabelString(option, label, 'labelKey');
    return Object.keys(label).length > 0 ? [[value, label]] : [];
  });
  return entries.length > 0 ? Object.fromEntries(entries) : undefined;
}

function parseFieldType(value: unknown): ConfigFieldType | undefined {
  return value === 'string' ||
    value === 'integer' ||
    value === 'number' ||
    value === 'boolean' ||
    value === 'object' ||
    value === 'array'
    ? value
    : undefined;
}

function assignNumber(raw: JsonRecord, property: ConfigSchemaProperty, key: keyof ConfigSchemaProperty) {
  if (typeof raw[key] === 'number') {
    property[key] = raw[key] as never;
  }
}

function assignString(raw: JsonRecord, property: ConfigSchemaProperty, key: keyof ConfigSchemaProperty) {
  if (typeof raw[key] === 'string') {
    property[key] = raw[key] as never;
  }
}

function assignOptionLabelString(
  raw: JsonRecord,
  property: ConfigSchemaOptionLabel,
  key: keyof ConfigSchemaOptionLabel,
) {
  if (typeof raw[key] === 'string') {
    property[key] = raw[key];
  }
}

function assignI18nString(raw: JsonRecord, property: ConfigSchemaI18n, key: keyof ConfigSchemaI18n) {
  if (typeof raw[key] === 'string') {
    property[key] = raw[key];
  }
}

function assignLegacyI18nString(
  raw: JsonRecord,
  property: ConfigSchemaI18n,
  legacyKey: string,
  key: keyof ConfigSchemaI18n,
) {
  if (typeof raw[legacyKey] === 'string' && !property[key]) {
    property[key] = raw[legacyKey];
  }
}
