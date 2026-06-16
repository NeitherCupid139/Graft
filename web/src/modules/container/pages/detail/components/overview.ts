type ContainerOverviewTagTheme = 'default' | 'primary' | 'success' | 'warning' | 'danger';

export type ContainerOverviewInfoRow =
  | {
      code?: boolean;
      copyValue?: string;
      displayValue: string;
      key: string;
      label: string;
      testId?: string;
      type: 'copy';
    }
  | {
      displayValue: string;
      key: string;
      label: string;
      type: 'text';
    }
  | {
      key: string;
      label: string;
      tagLabel: string;
      tagTheme: ContainerOverviewTagTheme;
      type: 'tag';
    }
  | {
      emptyLabel: string;
      key: string;
      label: string;
      ports: string[];
      type: 'ports';
    };

export type ContainerOverviewInfoSection = {
  key: string;
  rows: ContainerOverviewInfoRow[];
  title: string;
};
