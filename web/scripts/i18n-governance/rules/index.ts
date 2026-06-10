// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import type { I18nGovernanceRule } from '../types';
import { runLegacyRule } from './legacy-rule';
import { noFallbackOnlyKeyFirstRule } from './no-fallback-only-key-first';
import { noHardcodedPluginMessageRule } from './no-hardcoded-plugin-message';
import { noHardcodedUiPropRule } from './no-hardcoded-ui-prop';

const legacyRule: I18nGovernanceRule = {
  id: 'legacy',
  description:
    'Compatibility wrapper for existing locale catalog, missing key, template text, and unsafe datetime checks.',
  defaultSeverity: 'error',
  appliesTo: ['vue', 'ts', 'tsx', 'locale', 'go', 'schema'],
  check() {
    return runLegacyRule();
  },
};

export const rules: I18nGovernanceRule[] = [
  legacyRule,
  noHardcodedUiPropRule,
  noHardcodedPluginMessageRule,
  noFallbackOnlyKeyFirstRule,
];
