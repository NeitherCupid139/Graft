// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

export const ANNOUNCEMENT_ROUTE_PATH = {
  MANAGEMENT: '/server/announcements',
} as const;

export const ANNOUNCEMENT_API_PATH = {
  LIST: '/api/announcements',
  DETAIL: '/api/announcements/{id}',
  PUBLISH: '/api/announcements/{id}/publish',
  ARCHIVE: '/api/announcements/{id}/archive',
} as const;

export function buildAnnouncementDetailApiPath(id: number) {
  return ANNOUNCEMENT_API_PATH.DETAIL.replace('{id}', encodeURIComponent(String(id)));
}

export function buildAnnouncementPublishApiPath(id: number) {
  return ANNOUNCEMENT_API_PATH.PUBLISH.replace('{id}', encodeURIComponent(String(id)));
}

export function buildAnnouncementArchiveApiPath(id: number) {
  return ANNOUNCEMENT_API_PATH.ARCHIVE.replace('{id}', encodeURIComponent(String(id)));
}
