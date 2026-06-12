// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import type { paths } from '@/contracts/openapi/generated/schema';
import { request } from '@/utils/request';

import {
  ANNOUNCEMENT_API_PATH,
  buildAnnouncementArchiveApiPath,
  buildAnnouncementDetailApiPath,
  buildAnnouncementPublishApiPath,
} from '../contract/paths';
import type {
  AnnouncementItem,
  AnnouncementListQuery,
  AnnouncementListResponse,
  CreateAnnouncementRequest,
  PublishAnnouncementRequest,
  UpdateAnnouncementRequest,
} from '../types/announcement';

type AnnouncementListPath = (typeof ANNOUNCEMENT_API_PATH)['LIST'];
type GetAnnouncementsOperation = paths[AnnouncementListPath]['get'];
type GetAnnouncementsEnvelope = GetAnnouncementsOperation['responses'][200]['content']['application/json'];
type GetAnnouncementsData = NonNullable<GetAnnouncementsEnvelope['data']>;
type GetAnnouncementsQuery = NonNullable<GetAnnouncementsOperation['parameters']['query']>;

type PostAnnouncementsOperation = paths[AnnouncementListPath]['post'];
type PostAnnouncementsEnvelope = PostAnnouncementsOperation['responses'][201]['content']['application/json'];
type PostAnnouncementsData = NonNullable<PostAnnouncementsEnvelope['data']>;
type PostAnnouncementsBody = PostAnnouncementsOperation['requestBody']['content']['application/json'];

type AnnouncementDetailPath = (typeof ANNOUNCEMENT_API_PATH)['DETAIL'];
type GetAnnouncementOperation = paths[AnnouncementDetailPath]['get'];
type GetAnnouncementEnvelope = GetAnnouncementOperation['responses'][200]['content']['application/json'];
type GetAnnouncementData = NonNullable<GetAnnouncementEnvelope['data']>;
type GetAnnouncementPathParams = GetAnnouncementOperation['parameters']['path'];

type PutAnnouncementOperation = paths[AnnouncementDetailPath]['put'];
type PutAnnouncementEnvelope = PutAnnouncementOperation['responses'][200]['content']['application/json'];
type PutAnnouncementData = NonNullable<PutAnnouncementEnvelope['data']>;
type PutAnnouncementPathParams = PutAnnouncementOperation['parameters']['path'];
type PutAnnouncementBody = PutAnnouncementOperation['requestBody']['content']['application/json'];

type DeleteAnnouncementOperation = paths[AnnouncementDetailPath]['delete'];
type DeleteAnnouncementPathParams = DeleteAnnouncementOperation['parameters']['path'];

type AnnouncementPublishPath = (typeof ANNOUNCEMENT_API_PATH)['PUBLISH'];
type PostAnnouncementPublishOperation = paths[AnnouncementPublishPath]['post'];
type PostAnnouncementPublishEnvelope =
  PostAnnouncementPublishOperation['responses'][200]['content']['application/json'];
type PostAnnouncementPublishData = NonNullable<PostAnnouncementPublishEnvelope['data']>;
type PostAnnouncementPublishPathParams = PostAnnouncementPublishOperation['parameters']['path'];
type PostAnnouncementPublishBody = NonNullable<
  PostAnnouncementPublishOperation['requestBody']
>['content']['application/json'];

type AnnouncementArchivePath = (typeof ANNOUNCEMENT_API_PATH)['ARCHIVE'];
type PostAnnouncementArchiveOperation = paths[AnnouncementArchivePath]['post'];
type PostAnnouncementArchiveEnvelope =
  PostAnnouncementArchiveOperation['responses'][200]['content']['application/json'];
type PostAnnouncementArchiveData = NonNullable<PostAnnouncementArchiveEnvelope['data']>;
type PostAnnouncementArchivePathParams = PostAnnouncementArchiveOperation['parameters']['path'];

export function getAnnouncements(query?: AnnouncementListQuery) {
  return request.get<GetAnnouncementsData>({
    url: ANNOUNCEMENT_API_PATH.LIST,
    params: normalizeAnnouncementListQuery(query),
  }) as Promise<AnnouncementListResponse>;
}

export function createAnnouncement(payload: CreateAnnouncementRequest) {
  return request.post<PostAnnouncementsData>({
    url: ANNOUNCEMENT_API_PATH.LIST,
    data: payload as PostAnnouncementsBody,
  }) as Promise<AnnouncementItem>;
}

export function getAnnouncement(id: GetAnnouncementPathParams['id']) {
  return request.get<GetAnnouncementData>({
    url: buildAnnouncementDetailApiPath(id),
  }) as Promise<AnnouncementItem>;
}

export function updateAnnouncement(id: PutAnnouncementPathParams['id'], payload: UpdateAnnouncementRequest) {
  return request.put<PutAnnouncementData>({
    url: buildAnnouncementDetailApiPath(id),
    data: payload as PutAnnouncementBody,
  }) as Promise<AnnouncementItem>;
}

export function publishAnnouncement(id: PostAnnouncementPublishPathParams['id'], payload?: PublishAnnouncementRequest) {
  return request.post<PostAnnouncementPublishData>({
    url: buildAnnouncementPublishApiPath(id),
    data: payload as PostAnnouncementPublishBody | undefined,
  }) as Promise<AnnouncementItem>;
}

export function archiveAnnouncement(id: PostAnnouncementArchivePathParams['id']) {
  return request.post<PostAnnouncementArchiveData>({
    url: buildAnnouncementArchiveApiPath(id),
  }) as Promise<AnnouncementItem>;
}

export function deleteAnnouncement(id: DeleteAnnouncementPathParams['id']) {
  return request.delete<Record<string, never>>({
    url: buildAnnouncementDetailApiPath(id),
  });
}

export function normalizeAnnouncementListQuery(query?: AnnouncementListQuery): GetAnnouncementsQuery | undefined {
  if (!query) {
    return undefined;
  }

  return {
    ...(query.keyword ? { keyword: query.keyword } : {}),
    ...(query.level ? { level: query.level } : {}),
    ...(query.page ? { page: query.page } : {}),
    ...(query.page_size ? { page_size: query.page_size } : {}),
    ...(typeof query.pinned === 'boolean' ? { pinned: query.pinned } : {}),
    ...(query.sort ? { sort: query.sort } : {}),
    ...(query.status ? { status: query.status } : {}),
  } satisfies GetAnnouncementsQuery;
}
