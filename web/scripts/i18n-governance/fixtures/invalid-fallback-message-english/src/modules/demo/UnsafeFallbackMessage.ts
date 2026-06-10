// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

export function toNotification(error: Error) {
  return {
    messageKey: 'demo.request.failed',
    fallbackMessage: error.message || 'Request failed',
  };
}
