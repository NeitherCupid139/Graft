// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import { ANNOUNCEMENT_ROUTE_PATH } from './paths';

export const ANNOUNCEMENT_BOOTSTRAP_ROUTE = {
  MANAGEMENT: {
    menuPath: ANNOUNCEMENT_ROUTE_PATH.MANAGEMENT,
    routeName: 'AnnouncementManagement',
  },
} as const;

export type AnnouncementBootstrapRouteName =
  (typeof ANNOUNCEMENT_BOOTSTRAP_ROUTE)[keyof typeof ANNOUNCEMENT_BOOTSTRAP_ROUTE]['routeName'];
