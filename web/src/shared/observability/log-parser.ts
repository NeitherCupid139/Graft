// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import type { LogLevel, LogToken } from './log-highlight';
import { detectLogLevel, getLogLevelTone, normalizeLogLevel, tokenizeLogLine } from './log-highlight';

export type ParsedLogMetadata = Record<string, unknown>;
export type ParsedLogLine = {
  lineNo: number;
  timestamp: string;
  level: LogLevel | null;
  source: string;
  sourceShort: string;
  message: string;
  metadata: ParsedLogMetadata | null;
  raw: string;
  tone: ReturnType<typeof getLogLevelTone>;
};
export type DisplayLogLine = ParsedLogLine & {
  messageTokens: LogToken[];
  rawTokens: LogToken[];
  searchMatchCount: number;
};

const TIMESTAMP_PATTERN = /^(\d{4}-\d{2}-\d{2}(?:[T\s]\d{2}:\d{2}:\d{2}(?:[.,]\d+)?(?:Z|[+-]\d{2}:?\d{2})?)?)/;
const SOURCE_PATTERN = /^(?:[\w.-]+\/)*[\w.-]+\.(?:go|ts|tsx|vue|js|jsx|mjs|cjs|py|rs|java|kt|php|rb):\d+$/;
const KEY_VALUE_SUFFIX_PATTERN = /(?:^|\s)([A-Za-z_][\w.-]*)=("[^"]*"|'[^']*'|\S+)/g;
const METADATA_PRIORITY = ['request_id', 'status', 'duration', 'path', 'method', 'component', 'service', 'env'];
const LOW_SIGNAL_METADATA_PATTERNS = [/^legacy_/i, /^service$/i, /^env$/i];

export function parseLogLine(raw: string, lineNo: number): ParsedLogLine {
  const original = raw ?? '';
  const metadataResult = extractTrailingMetadata(original);
  const body = metadataResult.body.trim();
  const parsedHead = parseLogHead(body);
  const level = parsedHead.level ?? detectLogLevel(original);

  return {
    lineNo,
    timestamp: parsedHead.timestamp,
    level,
    source: parsedHead.source,
    sourceShort: shortenLogSource(parsedHead.source),
    message: parsedHead.message || body || original,
    metadata: metadataResult.metadata,
    raw: original,
    tone: getLogLevelTone(level),
  };
}

export function parseLogLines(lines: string[]): ParsedLogLine[] {
  return lines.map((line, index) => parseLogLine(line, index + 1));
}

export function buildDisplayLogLine(line: ParsedLogLine, keyword = ''): DisplayLogLine {
  return {
    ...line,
    messageTokens: tokenizeLogLine(line.message, keyword),
    rawTokens: tokenizeLogLine(line.raw, keyword),
    searchMatchCount: countKeywordMatches(line.raw, keyword),
  };
}

export function summarizeMetadata(metadata: ParsedLogMetadata | null, maxVisible = 3) {
  if (!metadata) {
    return { hiddenCount: 0, tags: [] as Array<[string, unknown]> };
  }

  const entries = Object.entries(metadata);
  const highSignalEntries = entries.filter(([key, value]) => !isLowSignalMetadata(key, value));
  const prioritized = [
    ...METADATA_PRIORITY.filter((key) => highSignalEntries.some(([entryKey]) => entryKey === key)).map(
      (key) => [key, metadata[key]] as [string, unknown],
    ),
    ...highSignalEntries.filter(([key]) => !METADATA_PRIORITY.includes(key)),
  ];

  return {
    hiddenCount: Math.max(0, entries.length - Math.min(prioritized.length, maxVisible)),
    tags: prioritized.slice(0, maxVisible),
  };
}

export function formatLogMetadataValue(value: unknown) {
  if (value === null) return 'null';
  if (value === undefined) return '';
  if (typeof value === 'string') return value;
  if (typeof value === 'number' || typeof value === 'boolean' || typeof value === 'bigint') return String(value);
  try {
    return JSON.stringify(value);
  } catch {
    return String(value);
  }
}

function countKeywordMatches(text: string, keyword = '') {
  const normalizedKeyword = keyword.trim().toLowerCase();
  if (!normalizedKeyword) {
    return 0;
  }

  const normalizedText = text.toLowerCase();
  let count = 0;
  let cursor = normalizedText.indexOf(normalizedKeyword);
  while (cursor >= 0) {
    count += 1;
    cursor = normalizedText.indexOf(normalizedKeyword, cursor + normalizedKeyword.length);
  }
  return count;
}

function parseLogHead(rawBody: string) {
  let body = rawBody.trim();
  let timestamp = '';
  let level: LogLevel | null = null;
  let source = '';

  const timestampMatch = TIMESTAMP_PATTERN.exec(body);
  if (timestampMatch) {
    timestamp = timestampMatch[1];
    body = body.slice(timestamp.length).trimStart();
  }

  const levelMatch = /^(FATAL|ERROR|WARN|WARNING|INFO|DEBUG|TRACE)\b/i.exec(body);
  if (levelMatch) {
    level = normalizeLogLevel(levelMatch[1]);
    body = body.slice(levelMatch[0].length).trimStart();
  }

  const firstToken = readToken(body);
  const secondToken = readToken(body.slice(firstToken.length).trimStart());
  if (SOURCE_PATTERN.test(firstToken)) {
    source = firstToken;
    body = body.slice(firstToken.length).trimStart();
  } else if (firstToken && SOURCE_PATTERN.test(secondToken)) {
    source = `${firstToken} ${secondToken}`;
    body = body.slice(firstToken.length).trimStart().slice(secondToken.length).trimStart();
  }

  return {
    level,
    message: body,
    source,
    timestamp,
  };
}

function extractTrailingMetadata(raw: string): { body: string; metadata: ParsedLogMetadata | null } {
  const jsonStart = findTrailingJsonStart(raw);
  if (jsonStart >= 0) {
    const jsonText = raw.slice(jsonStart).trim();
    try {
      const parsed = JSON.parse(jsonText);
      if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
        return {
          body: raw.slice(0, jsonStart).trimEnd(),
          metadata: parsed as ParsedLogMetadata,
        };
      }
    } catch {
      return { body: raw, metadata: null };
    }
  }

  return extractTrailingKeyValueMetadata(raw);
}

function extractTrailingKeyValueMetadata(raw: string): { body: string; metadata: ParsedLogMetadata | null } {
  const matches = [...raw.matchAll(KEY_VALUE_SUFFIX_PATTERN)];
  if (!matches.length) {
    return { body: raw, metadata: null };
  }

  const lastMatches: RegExpMatchArray[] = [];
  let expectedEnd = raw.length;
  for (let index = matches.length - 1; index >= 0; index -= 1) {
    const match = matches[index];
    const start = match.index ?? 0;
    const end = start + match[0].length;
    if (raw.slice(end, expectedEnd).trim() !== '') {
      break;
    }
    lastMatches.unshift(match);
    expectedEnd = start;
  }

  if (!lastMatches.length) {
    return { body: raw, metadata: null };
  }

  const metadata: ParsedLogMetadata = {};
  for (const match of lastMatches) {
    metadata[match[1]] = stripQuotes(match[2]);
  }

  return {
    body: raw.slice(0, lastMatches[0].index ?? 0).trimEnd(),
    metadata,
  };
}

function findTrailingJsonStart(raw: string) {
  let cursor = raw.lastIndexOf('{');
  while (cursor >= 0) {
    const suffix = raw.slice(cursor).trim();
    try {
      JSON.parse(suffix);
      return cursor;
    } catch {
      cursor = raw.lastIndexOf('{', cursor - 1);
    }
  }

  return -1;
}

function shortenLogSource(source: string) {
  if (!source) return '';
  const parts = source.split(/\s+/);
  const lastPart = parts.at(-1) ?? source;
  const pathParts = lastPart.split('/');
  const basename = pathParts.at(-1) ?? lastPart;
  if (parts.length > 1) {
    return basename;
  }
  if (pathParts.length >= 2 && basename === 'logger.go:61') {
    return `${pathParts.at(-2)}/${basename}`;
  }
  return basename;
}

function isLowSignalMetadata(key: string, value: unknown) {
  if (key === 'request_id' || key === 'status' || key === 'duration' || key === 'path' || key === 'method') {
    return false;
  }
  if (key === 'component') {
    return false;
  }
  if (key === 'service' && typeof value === 'string') {
    return value === 'sub2api' || value === 'sub2api-api';
  }
  if (key === 'env' && typeof value === 'string') {
    return value === 'production';
  }
  return LOW_SIGNAL_METADATA_PATTERNS.some((pattern) => pattern.test(key));
}

function readToken(value: string) {
  return value.trimStart().split(/\s+/, 1)[0] ?? '';
}

function stripQuotes(value: string) {
  return value.replace(/^["']|["']$/g, '');
}
