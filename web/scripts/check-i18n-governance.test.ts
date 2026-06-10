// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import { spawnSync } from 'node:child_process';
import { cpSync, mkdirSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';

import { afterEach, describe, expect, it } from 'vitest';

const tempRoots: string[] = [];

function createTempWebRoot(source: string) {
  const repoRoot = join(tmpdir(), `graft-i18n-governance-${process.pid}-${tempRoots.length}`);
  const root = join(repoRoot, 'web');
  tempRoots.push(repoRoot);

  mkdirSync(join(root, 'scripts'), { recursive: true });
  mkdirSync(join(root, 'src/modules/demo/locales'), { recursive: true });
  mkdirSync(join(repoRoot, 'server/internal'), { recursive: true });
  mkdirSync(join(repoRoot, 'server/modules'), { recursive: true });
  cpSync(join(process.cwd(), 'scripts/check-i18n-governance.ts'), join(root, 'scripts/check-i18n-governance.ts'));
  writeFileSync(join(root, 'src/modules/demo/UnsafeTime.vue'), source);
  writeFileSync(join(root, 'src/modules/demo/locales/en-US.json'), '{}');
  writeFileSync(join(root, 'src/modules/demo/locales/zh-CN.json'), '{}');

  return root;
}

async function runGovernanceScript(source: string) {
  const root = createTempWebRoot(source);
  const process = spawnSync('bun', ['run', 'scripts/check-i18n-governance.ts'], {
    cwd: root,
    encoding: 'utf8',
  });

  return {
    stdout: process.stdout,
    stderr: process.stderr,
    exitCode: process.status,
  };
}

afterEach(() => {
  for (const root of tempRoots.splice(0)) {
    rmSync(root, { force: true, recursive: true });
  }
});

describe('check-i18n-governance datetime formatting scan', () => {
  it('blocks visible datetime formatting that depends on the host locale', async () => {
    const result = await runGovernanceScript(`
<template>{{ label }}</template>
<script setup lang="ts">
const label = new Intl.DateTimeFormat(undefined, { dateStyle: 'medium' }).format(new Date());
const fallback = new Date().toLocaleString();
</script>
`);

    expect(result.exitCode).toBe(1);
    expect(result.stdout).toContain('visible datetime formatting must pass the active locale instead of undefined');
    expect(result.stdout).toContain('toLocaleString must pass the active locale');
    expect(result.stderr).toBe('');
  });

  it('allows visible datetime formatting through an explicit locale-aware formatter', async () => {
    const result = await runGovernanceScript(`
<template>{{ label }}</template>
<script setup lang="ts">
import { formatLocaleDateTime } from '@/shared/observability';
const label = formatLocaleDateTime('2026-06-10T02:38:00Z', 'en-US');
</script>
`);

    expect(result.exitCode).toBe(0);
    expect(result.stdout).toContain('No hard-coded UI text or locale governance issues found.');
    expect(result.stderr).toBe('');
  });
});
