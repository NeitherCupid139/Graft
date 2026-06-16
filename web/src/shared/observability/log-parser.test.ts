// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import { describe, expect, it } from 'vitest';

import { buildDisplayLogLine, type ParsedLogMetadata, parseLogLine, summarizeMetadata } from './log-parser';

describe('log-parser', () => {
  it('parses ordinary INFO logs into structured fields', () => {
    const line = parseLogLine(
      '2026-06-17T06:31:40.491+0800 INFO service/pricing_service.go:145 [Pricing] Service initialized',
      8,
    );

    expect(line.lineNo).toBe(8);
    expect(line.timestamp).toBe('2026-06-17T06:31:40.491+0800');
    expect(line.level).toBe('INFO');
    expect(line.source).toBe('service/pricing_service.go:145');
    expect(line.sourceShort).toBe('pricing_service.go:145');
    expect(line.message).toBe('[Pricing] Service initialized');
    expect(line.metadata).toBeNull();
  });

  it('parses WARN logs and maps warning tone without strong row coloring', () => {
    const line = parseLogLine(
      '2026-06-17T06:31:40.497+0800 WARN stdlog server/wire_gen.go:262 Warning: server.trusted_proxies is empty',
      15,
    );

    expect(line.level).toBe('WARN');
    expect(line.source).toBe('stdlog server/wire_gen.go:262');
    expect(line.sourceShort).toBe('wire_gen.go:262');
    expect(line.tone).toBe('warning');
    expect(line.message).toBe('Warning: server.trusted_proxies is empty');
  });

  it('extracts trailing JSON metadata and summarizes common fields first', () => {
    const line = parseLogLine(
      '2026-06-17T06:31:42.585+0800 INFO middleware/logger.go:61 http request completed {"service":"sub2api","env":"production","component":"http","request_id":"abc","duration":"12ms","method":"GET","path":"/health","status":200}',
      29,
    );
    const metadata = line.metadata as ParsedLogMetadata;
    const summary = summarizeMetadata(metadata, 3);

    expect(line.message).toBe('http request completed');
    expect(metadata.service).toBe('sub2api');
    expect(metadata.status).toBe(200);
    expect(summary.tags.map(([key]) => key)).toEqual(['request_id', 'status', 'duration']);
    expect(summary.hiddenCount).toBe(5);
  });

  it('folds repeated low-signal metadata out of the default tag summary', () => {
    const line = parseLogLine(
      '2026-06-17T06:43:20.498+0800 INFO stdlog service/timing_wheel.go:37 flushed {"service":"sub2api","env":"production","legacy_stdlog":true,"request_id":"req-1","component":"scheduler"}',
      171,
    );
    const summary = summarizeMetadata(line.metadata, 3);

    expect(summary.tags.map(([key]) => key)).toEqual(['request_id', 'component']);
    expect(summary.hiddenCount).toBe(3);
  });

  it('falls back to raw text when trailing JSON parsing fails', () => {
    const raw =
      '2026-06-17T06:31:40.491+0800 INFO service/pricing_service.go:145 [Pricing] Service initialized {"service":';
    const line = parseLogLine(raw, 1);

    expect(line.metadata).toBeNull();
    expect(line.raw).toBe(raw);
    expect(line.message).toContain('{"service":');
  });

  it('highlights search keywords without changing log order', () => {
    const line = buildDisplayLogLine(
      parseLogLine(
        '2026-06-17T06:31:42.585+0800 INFO middleware/logger.go:61 http request completed {"request_id":"abc"}',
        2,
      ),
      'request',
    );

    expect(line.lineNo).toBe(2);
    expect(line.searchMatchCount).toBe(2);
    expect(line.messageTokens.some((token) => token.type === 'keyword' && token.text === 'request')).toBe(true);
  });
});
