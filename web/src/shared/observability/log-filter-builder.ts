export type LogFilterOption = {
  label: string;
  value: string;
};

export type LogFilterFieldKind = 'text' | 'select' | 'multi-select' | 'tag-input' | 'special';

export type LogFilterFieldDefinition = {
  key: string;
  label: string;
  kind: LogFilterFieldKind;
  placeholder?: string;
  options?: LogFilterOption[];
  disabled?: boolean;
};

export type LogFilterTag = {
  key: string;
  label: string;
  closable?: boolean;
};

export type LogTimeRangeField = {
  key: string;
  label: string;
  value: string[];
  placeholder: [string, string];
};

export type LogSortItem = {
  field: string;
  direction?: 'asc' | 'desc';
};

export type LogSortOption = {
  label: string;
  value: string;
};

export type LogFilterPreset = {
  key: string;
  title: string;
};

export type LogSorterUiState = {
  sortAddDisabled: boolean;
  sortFieldOptionsByIndex: LogSortOption[][];
  sortMoveDownDisabled: boolean[];
  sortMoveUpDisabled: boolean[];
  sorters: LogSortItem[];
};
