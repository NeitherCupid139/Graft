import type { RealtimeSubscriptionResponse } from './api';
import { postRealtimeSubscription } from './api';
import { toRealtimeWebSocketUrl } from './url';

type RealtimeSocketState = 'idle' | 'connecting' | 'open' | 'closed' | 'error';

type OpenRealtimeTopicSocketOptions<TMessage> = {
  topic: string;
  issueTicket?: (topic: string) => Promise<RealtimeSubscriptionResponse>;
  onMessage?: (message: TMessage) => void;
  onStateChange?: (state: RealtimeSocketState) => void;
  onError?: (message: string) => void;
  parseMessage?: (raw: unknown) => TMessage | null;
};

export type RealtimeTopicSocketController = {
  close: () => void;
  reconnect: () => void;
};

function defaultParseMessage<TMessage>(raw: unknown) {
  if (typeof raw !== 'string') {
    return null;
  }
  try {
    return JSON.parse(raw) as TMessage;
  } catch {
    return null;
  }
}

export function openRealtimeTopicSocket<TMessage>(
  options: OpenRealtimeTopicSocketOptions<TMessage>,
): RealtimeTopicSocketController {
  let socket: WebSocket | null = null;
  let reconnectTimer: number | null = null;
  let closed = false;
  let connectionId = 0;

  function clearReconnectTimer() {
    if (reconnectTimer !== null) {
      clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }
  }

  function emitState(nextState: RealtimeSocketState) {
    options.onStateChange?.(nextState);
  }

  async function connect() {
    clearReconnectTimer();
    if (closed) {
      return;
    }

    const currentConnectionId = ++connectionId;
    emitState('connecting');
    try {
      const issued = await (options.issueTicket?.(options.topic) ?? postRealtimeSubscription({ topic: options.topic }));
      if (closed || currentConnectionId !== connectionId) {
        return;
      }
      const nextSocket = new WebSocket(toRealtimeWebSocketUrl(issued.websocket_url));
      socket = nextSocket;

      nextSocket.onopen = () => {
        if (socket !== nextSocket || closed || currentConnectionId !== connectionId) {
          return;
        }
        emitState('open');
      };

      nextSocket.onmessage = (event) => {
        if (socket !== nextSocket || closed || currentConnectionId !== connectionId) {
          return;
        }
        const parser = options.parseMessage ?? defaultParseMessage<TMessage>;
        const parsed = parser(event.data);
        if (!parsed) {
          return;
        }
        options.onMessage?.(parsed);
      };

      nextSocket.onerror = () => {
        if (socket !== nextSocket || closed || currentConnectionId !== connectionId) {
          return;
        }
        emitState('error');
        options.onError?.('WebSocket transport error');
      };

      nextSocket.onclose = () => {
        if (socket !== nextSocket || currentConnectionId !== connectionId) {
          return;
        }
        socket = null;
        if (closed) {
          emitState('idle');
          return;
        }
        emitState('closed');
        reconnectTimer = window.setTimeout(() => {
          void connect();
        }, 1000);
      };
    } catch (error) {
      if (closed || currentConnectionId !== connectionId) {
        return;
      }
      emitState('error');
      options.onError?.(error instanceof Error ? error.message : 'Failed to issue realtime subscription ticket');
      reconnectTimer = window.setTimeout(() => {
        void connect();
      }, 1000);
    }
  }

  function close() {
    closed = true;
    clearReconnectTimer();
    if (socket) {
      socket.onopen = null;
      socket.onmessage = null;
      socket.onerror = null;
      socket.onclose = null;
      if (socket.readyState < WebSocket.CLOSING) {
        socket.close();
      }
    }
    socket = null;
    emitState('idle');
  }

  function reconnect() {
    closed = false;
    connectionId += 1;
    if (socket && socket.readyState < WebSocket.CLOSING) {
      socket.close();
    }
    socket = null;
    void connect();
  }

  void connect();

  return {
    close,
    reconnect,
  };
}
