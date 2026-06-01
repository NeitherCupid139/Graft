export type SortDirection = 'asc' | 'desc';

export type QuerySorter<Field extends string = string> = {
  field: Field;
  direction?: SortDirection;
};

export type SorterState<Field extends string, TState extends { sorters: QuerySorter<Field>[] }> = TState;

type SortOption<Field extends string> = {
  label: string;
  value: Field;
};

export function createSingleSorter<Field extends string>(
  field?: Field | null,
  direction: SortDirection = 'desc',
): QuerySorter<Field>[] {
  if (!field) {
    return [];
  }

  return [{ field, direction }];
}

function createQuerySorter<Field extends string>(
  field?: Field | null,
  direction: SortDirection = 'desc',
): QuerySorter<Field> | null {
  if (!field) {
    return null;
  }

  return { field, direction };
}

export function appendSorter<Field extends string>(
  sorters: QuerySorter<Field>[],
  field?: Field | null,
  direction: SortDirection = 'desc',
) {
  const sorter = createQuerySorter(field, direction);
  if (!sorter) {
    return sorters;
  }

  return [...sorters, sorter];
}

function replaceSorterField<Field extends string>(
  sorters: QuerySorter<Field>[],
  index: number,
  field: Field | '',
  fallbackDirection: SortDirection = 'desc',
) {
  const nextSorters = [...sorters];
  if (!field) {
    nextSorters.splice(index, 1);
    return nextSorters;
  }

  nextSorters[index] = {
    field,
    direction: nextSorters[index]?.direction ?? fallbackDirection,
  };
  return nextSorters;
}

function replaceSorterDirection<Field extends string>(
  sorters: QuerySorter<Field>[],
  index: number,
  direction: SortDirection,
) {
  const sorter = sorters[index];
  if (!sorter?.field) {
    return sorters;
  }

  const nextSorters = [...sorters];
  nextSorters[index] = {
    field: sorter.field,
    direction,
  };
  return nextSorters;
}

function replaceSorterFieldFromInput<Field extends string>(
  sorters: QuerySorter<Field>[],
  index: number,
  value: string | number | Array<string | number> | undefined,
  normalizeField: (value: string) => Field | '',
  fallbackDirection: SortDirection = 'desc',
) {
  const field = typeof value === 'string' ? normalizeField(value) : '';
  return replaceSorterField(sorters, index, field, fallbackDirection);
}

function replaceSorterDirectionFromInput<Field extends string>(
  sorters: QuerySorter<Field>[],
  index: number,
  value: string | number | Array<string | number> | undefined,
  normalizeDirection: (value: string) => SortDirection,
) {
  return replaceSorterDirection(sorters, index, normalizeDirection(typeof value === 'string' ? value : ''));
}

export function withUpdatedSorters<Field extends string, TState extends { sorters: QuerySorter<Field>[] }>(
  state: TState,
  sorters: QuerySorter<Field>[],
): TState {
  return {
    ...state,
    sorters,
  };
}

export function withSorterFieldFromInput<Field extends string, TState extends { sorters: QuerySorter<Field>[] }>(
  state: TState,
  index: number,
  value: string | number | Array<string | number> | undefined,
  normalizeField: (value: string) => Field | '',
  fallbackDirection: SortDirection = 'desc',
): TState {
  return withUpdatedSorters(
    state,
    replaceSorterFieldFromInput(state.sorters, index, value, normalizeField, fallbackDirection),
  );
}

export function withSorterDirectionFromInput<Field extends string, TState extends { sorters: QuerySorter<Field>[] }>(
  state: TState,
  index: number,
  value: string | number | Array<string | number> | undefined,
  normalizeDirection: (value: string) => SortDirection,
): TState {
  return withUpdatedSorters(state, replaceSorterDirectionFromInput(state.sorters, index, value, normalizeDirection));
}

function buildSorterTagLabel<Field extends string>(
  sorter: QuerySorter<Field>,
  options: Array<SortOption<Field>>,
  prefix: string,
  index: number,
) {
  const fieldLabel = options.find((option) => option.value === sorter.field)?.label ?? sorter.field;
  const arrow = sorter.direction === 'asc' ? '↑' : sorter.direction === 'desc' ? '↓' : '';
  return `${prefix} ${index + 1}：${[fieldLabel, arrow].filter(Boolean).join(' ')}`;
}

export function prependSorterTags<Key extends string, Field extends string>(
  tags: Array<{ key: Key; label: string }>,
  sorters: QuerySorter<Field>[],
  options: Array<SortOption<Field>>,
  prefix: string,
): Array<{ key: Key | `sorter:${number}`; label: string }> {
  if (!sorters.length) {
    return tags;
  }

  return [
    ...sorters.map((sorter, index) => ({
      key: `sorter:${index}` as const,
      label: buildSorterTagLabel(sorter, options, prefix, index),
    })),
    ...tags,
  ];
}

export function encodeSorters<Field extends string>(sorters: QuerySorter<Field>[]) {
  return sorters
    .map((sorter) => {
      const field = String(sorter.field || '').trim();
      if (!field) {
        return '';
      }

      const direction = sorter.direction === 'asc' ? 'asc' : 'desc';
      return `${field}:${direction}`;
    })
    .filter(Boolean);
}

export function decodeSorters<Field extends string>(
  rawValue: string | string[] | undefined,
  normalizeField: (value: string) => Field | '',
  normalizeDirection: (value: string) => SortDirection,
) {
  const values = Array.isArray(rawValue) ? rawValue : rawValue ? [rawValue] : [];

  return values.reduce<QuerySorter<Field>[]>((acc, value) => {
    const candidate = value.trim();
    if (!candidate) {
      return acc;
    }

    const [rawField = '', rawDirection = 'desc'] = candidate.split(':');
    const field = normalizeField(rawField.trim());
    if (!field) {
      return acc;
    }

    acc.push({
      field,
      direction: normalizeDirection(rawDirection.trim() || 'desc'),
    });
    return acc;
  }, []);
}
