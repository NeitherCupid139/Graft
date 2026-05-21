export const USER_STATUS = {
  DISABLED: 'disabled',
  ENABLED: 'enabled',
} as const;

export type UserStatus = (typeof USER_STATUS)[keyof typeof USER_STATUS];
