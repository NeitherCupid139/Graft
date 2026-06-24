import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { openRealtimeTopicSocket } from './ws-client';

class MockWebSocket {
  static instances: MockWebSocket[] = [];
  static readonly CONNECTING = 0;
  static readonly OPEN = 1;
  static readonly CLOSING = 2;
  static readonly CLOSED = 3;

  readonly url: string;
  readyState = MockWebSocket.CONNECTING;
  onopen: (() => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;
  onerror: (() => void) | null = null;
  onclose: (() => void) | null = null;
  close = vi.fn(() => {
    this.readyState = MockWebSocket.CLOSING;
  });

  constructor(url: string) {
    this.url = url;
    MockWebSocket.instances.push(this);
  }

  emitOpen() {
    this.readyState = MockWebSocket.OPEN;
    this.onopen?.();
  }

  emitMessage(data: unknown) {
    this.onmessage?.({ data } as MessageEvent);
  }

  emitClose() {
    this.readyState = MockWebSocket.CLOSED;
    this.onclose?.();
  }
}

describe('openRealtimeTopicSocket', () => {
  const originalWebSocket = globalThis.WebSocket;

  beforeEach(() => {
    MockWebSocket.instances = [];
    vi.useFakeTimers();
    globalThis.WebSocket = Object.assign(MockWebSocket, {
      CONNECTING: MockWebSocket.CONNECTING,
      OPEN: MockWebSocket.OPEN,
      CLOSING: MockWebSocket.CLOSING,
      CLOSED: MockWebSocket.CLOSED,
    }) as unknown as typeof WebSocket;
  });

  afterEach(() => {
    vi.useRealTimers();
    if (originalWebSocket) {
      globalThis.WebSocket = originalWebSocket;
    } else {
      Reflect.deleteProperty(globalThis, 'WebSocket');
    }
  });

  it('issues a ticket before opening the unified websocket url', async () => {
    const issueTicket = vi.fn().mockResolvedValue({
      topic: 'container.stats:container-1',
      ticket: 'opaque-ticket',
      websocket_url: '/ws?topic=container.stats%3Acontainer-1&ticket=opaque-ticket',
      expires_at: '2026-06-24T08:00:30Z',
    });
    const onMessage = vi.fn();

    openRealtimeTopicSocket({
      topic: 'container.stats:container-1',
      issueTicket,
      onMessage,
      parseMessage: (raw) => (typeof raw === 'string' ? (JSON.parse(raw) as { value: number }) : null),
    });

    await vi.runAllTicks();

    expect(issueTicket).toHaveBeenCalledWith('container.stats:container-1');
    expect(MockWebSocket.instances).toHaveLength(1);
    expect(MockWebSocket.instances[0]?.url).toContain('/ws?topic=container.stats%3Acontainer-1&ticket=opaque-ticket');

    MockWebSocket.instances[0]?.emitOpen();
    MockWebSocket.instances[0]?.emitMessage(JSON.stringify({ value: 42 }));

    expect(onMessage).toHaveBeenCalledWith({ value: 42 });
  });
});
