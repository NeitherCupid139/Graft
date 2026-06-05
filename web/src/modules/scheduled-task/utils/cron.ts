export type CronValidationResult = {
  valid: boolean;
  messageKey?: CronValidationMessageKey;
  messageParams?: Record<string, string | number>;
};

export type CronValidationMessageKey =
  | 'scheduledTask.cronValidation.fieldCount'
  | 'scheduledTask.cronValidation.stepRange'
  | 'scheduledTask.cronValidation.fieldRange';

export type CronDescriptionKey =
  | 'scheduledTask.cronDescription.everyMinute'
  | 'scheduledTask.cronDescription.everyNMinutes'
  | 'scheduledTask.cronDescription.hourly'
  | 'scheduledTask.cronDescription.daily'
  | 'scheduledTask.cronDescription.weekly'
  | 'scheduledTask.cronDescription.monthly'
  | 'scheduledTask.cronDescription.custom'
  | 'scheduledTask.cronDescription.invalid';

export type CronDescriptionResult = {
  key: CronDescriptionKey;
  params?: Record<string, string | number>;
  normalizedExpression: string;
  valid: boolean;
};

export type CronDescriptionTranslate = (key: CronDescriptionKey, params?: Record<string, string | number>) => string;

const FIELD_COUNT_UNIX = 5;
const FIELD_COUNT_SECONDS = 6;

type CronFieldRule = {
  name: string;
  min: number;
  max: number;
  allowStep: boolean;
};

const CRON_FIELD_RULES: CronFieldRule[] = [
  { name: 'seconds', min: 0, max: 59, allowStep: true },
  { name: 'minutes', min: 0, max: 59, allowStep: true },
  { name: 'hours', min: 0, max: 23, allowStep: false },
  { name: 'day-of-month', min: 1, max: 31, allowStep: false },
  { name: 'month', min: 1, max: 12, allowStep: false },
  { name: 'day-of-week', min: 0, max: 7, allowStep: false },
];

export function normalizeCronExpression(expression: string): string {
  const normalized = splitCronFields(expression).join(' ');
  const fields = splitCronFields(normalized);

  if (fields.length === FIELD_COUNT_UNIX) {
    return ['0', ...fields].join(' ');
  }

  return normalized;
}

export function validateCronExpression(expression: string): CronValidationResult {
  const fields = splitCronFields(expression);

  if (fields.length !== FIELD_COUNT_UNIX && fields.length !== FIELD_COUNT_SECONDS) {
    return {
      valid: false,
      messageKey: 'scheduledTask.cronValidation.fieldCount',
      messageParams: { unixFields: FIELD_COUNT_UNIX, secondsFields: FIELD_COUNT_SECONDS },
    };
  }

  const normalizedFields = normalizeCronExpression(expression).split(' ');

  for (const [index, field] of normalizedFields.entries()) {
    const rule = CRON_FIELD_RULES[index];
    const result = validateCronField(field, rule);
    if (!result.valid) {
      return result;
    }
  }

  return { valid: true };
}

export function describeCronExpression(
  expression: string,
  translate?: CronDescriptionTranslate,
): CronDescriptionResult | string {
  const validation = validateCronExpression(expression);
  const normalizedExpression = normalizeCronExpression(expression);

  const description = validation.valid
    ? describeNormalizedCronExpression(normalizedExpression)
    : ({
        key: 'scheduledTask.cronDescription.invalid',
        params: { expression: normalizedExpression },
        normalizedExpression,
        valid: false,
      } satisfies CronDescriptionResult);

  if (translate) {
    return translate(description.key, description.params);
  }

  return description;
}

function describeNormalizedCronExpression(normalizedExpression: string): CronDescriptionResult {
  const fields = normalizedExpression.split(' ');
  const [second, minute, hour, dayOfMonth, month, dayOfWeek] = fields;

  if (second === '0' && minute === '*' && hour === '*' && dayOfMonth === '*' && month === '*' && dayOfWeek === '*') {
    return {
      key: 'scheduledTask.cronDescription.everyMinute',
      normalizedExpression,
      valid: true,
    };
  }

  const minuteInterval = parseStepValue(minute);
  if (
    second === '0' &&
    minuteInterval !== null &&
    hour === '*' &&
    dayOfMonth === '*' &&
    month === '*' &&
    dayOfWeek === '*'
  ) {
    return {
      key: 'scheduledTask.cronDescription.everyNMinutes',
      params: { interval: minuteInterval },
      normalizedExpression,
      valid: true,
    };
  }

  if (second === '0' && minute === '0' && hour === '*' && dayOfMonth === '*' && month === '*' && dayOfWeek === '*') {
    return {
      key: 'scheduledTask.cronDescription.hourly',
      normalizedExpression,
      valid: true,
    };
  }

  if (second === '0' && minute === '0' && isNumberInRange(hour, 0, 23) && dayOfMonth === '*' && month === '*') {
    if (dayOfWeek === '*') {
      return {
        key: 'scheduledTask.cronDescription.daily',
        params: { hour: Number(hour) },
        normalizedExpression,
        valid: true,
      };
    }

    if (isNumberInRange(dayOfWeek, 0, 7)) {
      return {
        key: 'scheduledTask.cronDescription.weekly',
        params: { hour: Number(hour), dayOfWeek: Number(dayOfWeek) },
        normalizedExpression,
        valid: true,
      };
    }
  }

  if (
    second === '0' &&
    minute === '0' &&
    isNumberInRange(hour, 0, 23) &&
    isNumberInRange(dayOfMonth, 1, 31) &&
    month === '*' &&
    dayOfWeek === '*'
  ) {
    return {
      key: 'scheduledTask.cronDescription.monthly',
      params: { hour: Number(hour), dayOfMonth: Number(dayOfMonth) },
      normalizedExpression,
      valid: true,
    };
  }

  return {
    key: 'scheduledTask.cronDescription.custom',
    params: { expression: normalizedExpression },
    normalizedExpression,
    valid: true,
  };
}

function splitCronFields(expression: string): string[] {
  return expression.trim().split(/\s+/).filter(Boolean);
}

function validateCronField(field: string, rule: CronFieldRule): CronValidationResult {
  if (field === '*') {
    return { valid: true };
  }

  if (rule.allowStep && field.startsWith('*/')) {
    const step = field.slice(2);
    if (isNumberInRange(step, 1, rule.max)) {
      return { valid: true };
    }

    return {
      valid: false,
      messageKey: 'scheduledTask.cronValidation.stepRange',
      messageParams: { field: rule.name, min: 1, max: rule.max },
    };
  }

  if (isNumberInRange(field, rule.min, rule.max)) {
    return { valid: true };
  }

  return {
    valid: false,
    messageKey: 'scheduledTask.cronValidation.fieldRange',
    messageParams: { field: rule.name, min: rule.min, max: rule.max },
  };
}

function parseStepValue(field: string): number | null {
  if (!field.startsWith('*/')) {
    return null;
  }

  const step = field.slice(2);
  return isPositiveInteger(step) ? Number(step) : null;
}

function isNumberInRange(value: string, min: number, max: number): boolean {
  if (!/^\d+$/.test(value)) {
    return false;
  }

  const numericValue = Number(value);
  return numericValue >= min && numericValue <= max;
}

function isPositiveInteger(value: string): boolean {
  return /^[1-9]\d*$/.test(value);
}
