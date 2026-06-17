// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import { nextTick, readonly, ref } from 'vue';

export const ROUTE_LOADING_MIN_MS = 150;
export const ROUTE_LOADING_MAX_MS = 5000;

const loading = ref(false);
let loadingStartedAt = 0;
let loadingToken = 0;
let minTimer: ReturnType<typeof setTimeout> | undefined;
let maxTimer: ReturnType<typeof setTimeout> | undefined;

export const routeLoading = readonly(loading);

function clearTimers() {
  if (minTimer) {
    clearTimeout(minTimer);
    minTimer = undefined;
  }

  if (maxTimer) {
    clearTimeout(maxTimer);
    maxTimer = undefined;
  }
}

function requestNextFrame() {
  return new Promise<void>((resolve) => {
    if (typeof requestAnimationFrame === 'function') {
      requestAnimationFrame(() => resolve());
      return;
    }

    setTimeout(resolve, 0);
  });
}

function stopRouteLoadingNow() {
  clearTimers();
  loading.value = false;
}

export function startRouteLoading() {
  loadingToken += 1;
  loadingStartedAt = Date.now();
  clearTimers();
  loading.value = true;
  maxTimer = setTimeout(stopRouteLoadingNow, ROUTE_LOADING_MAX_MS);
}

export async function finishRouteLoadingAfterRender() {
  const token = loadingToken;
  await nextTick();
  await requestNextFrame();

  if (token !== loadingToken) {
    return;
  }

  const remaining = ROUTE_LOADING_MIN_MS - (Date.now() - loadingStartedAt);
  if (remaining <= 0) {
    stopRouteLoadingNow();
    return;
  }

  minTimer = setTimeout(stopRouteLoadingNow, remaining);
}

export function hideRouteLoading() {
  loadingToken += 1;
  stopRouteLoadingNow();
}
