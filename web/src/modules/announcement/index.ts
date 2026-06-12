// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import type { WebModuleRegistration } from '@/modules/types';

import { announcementBootstrapRouteRegistrations } from './bootstrap-routes';
import { ANNOUNCEMENT_PERMISSION_CODE } from './contract/permissions';

export const announcementModuleRegistration: WebModuleRegistration = {
  moduleId: 'announcement',
  bootstrapRoutes: announcementBootstrapRouteRegistrations,
};

export const announcementModulePermissionCodes = ANNOUNCEMENT_PERMISSION_CODE;

export default announcementModuleRegistration;
