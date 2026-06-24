export function toRealtimeWebSocketUrl(relativePath: string) {
  const base = new URL(window.location.href);
  const protocol = base.protocol === 'https:' ? 'wss:' : 'ws:';
  return new URL(relativePath, `${protocol}//${base.host}`).toString();
}
