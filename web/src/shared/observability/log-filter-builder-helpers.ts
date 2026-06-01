import { computed, type Ref } from 'vue';

import { buildSorterUiState, type QuerySorter, type SortDirection, type SortOption } from './sorters';

type BuilderListenerConfig<TPreset extends string, TFieldKey extends string, TTimeValue> = {
  addSorter: () => void;
  clearTag: (key: string) => void;
  emitApplyPreset: (preset: TPreset) => void;
  emitReset: () => void;
  emitSearch: () => void;
  handleFieldUpdate: (payload: { key: string; value: string | string[] }) => void;
  moveSorterDown: (index: number) => void;
  moveSorterUp: (index: number) => void;
  removeSorter: (index: number) => void;
  selectedFieldKey: Ref<TFieldKey>;
  updateKeyword: (value: string) => void;
  updateSortDirection: (index: number, value: string | number | Array<string | number> | undefined) => void;
  updateSortField: (index: number, value: string | number | Array<string | number> | undefined) => void;
  updateTimeField: (payload: TTimeValue) => void;
};

export function useLogSorterUiState<Field extends string>(
  sorters: () => QuerySorter<Field>[],
  sortOptions: () => Array<SortOption<Field>>,
) {
  const sorterUiState = computed(() => buildSorterUiState(sorters(), sortOptions()));

  return {
    normalizedSorters: computed(() => sorterUiState.value.sorters),
    sortAddDisabled: computed(() => sorterUiState.value.sortAddDisabled),
    sortFieldOptionsByIndex: computed(() => sorterUiState.value.sortFieldOptionsByIndex),
    sortMoveDownDisabled: computed(() => sorterUiState.value.sortMoveDownDisabled),
    sortMoveUpDisabled: computed(() => sorterUiState.value.sortMoveUpDisabled),
  };
}

export function createLogBuilderListeners<TPreset extends string, TFieldKey extends string, TTimeValue>(
  config: BuilderListenerConfig<TPreset, TFieldKey, TTimeValue>,
) {
  return {
    'add-sorter': config.addSorter,
    'apply-preset': (preset: string) => config.emitApplyPreset(preset as TPreset),
    'close-tag': (key: string) => config.clearTag(key),
    'move-sorter-down': config.moveSorterDown,
    'move-sorter-up': config.moveSorterUp,
    'remove-sorter': config.removeSorter,
    reset: config.emitReset,
    search: config.emitSearch,
    'update:field': config.handleFieldUpdate,
    'update:keyword': (value: string) => config.updateKeyword(value),
    'update:selected-field-key': (key: string) => {
      config.selectedFieldKey.value = key as TFieldKey;
    },
    'update:sort-direction': ({
      index,
      value,
    }: {
      index: number;
      value: string | number | Array<string | number> | undefined;
    }) => config.updateSortDirection(index, value),
    'update:sort-field': ({
      index,
      value,
    }: {
      index: number;
      value: string | number | Array<string | number> | undefined;
    }) => config.updateSortField(index, value),
    'update:time-field': (payload: TTimeValue) => config.updateTimeField(payload),
  };
}

export function createSortDirection(value: string): SortDirection {
  return value === 'asc' ? 'asc' : 'desc';
}
