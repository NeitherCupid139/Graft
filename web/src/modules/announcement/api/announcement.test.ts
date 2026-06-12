// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import { describe, expect, it } from 'vitest';

import { normalizeAnnouncementListQuery } from './announcement';

describe('announcement API query mapping', () => {
  it('omits empty filters and preserves typed backend parameters', () => {
    expect(
      normalizeAnnouncementListQuery({
        keyword: '',
        level: undefined,
        page: 1,
        page_size: 20,
        pinned: false,
        sort: 'pinned_publish_desc',
        status: 'published',
      }),
    ).toEqual({
      page: 1,
      page_size: 20,
      pinned: false,
      sort: 'pinned_publish_desc',
      status: 'published',
    });
  });

  it('returns undefined for absent query objects', () => {
    expect(normalizeAnnouncementListQuery()).toBeUndefined();
  });
});
