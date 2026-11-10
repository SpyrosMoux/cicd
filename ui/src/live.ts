// The collector sends text only. A frame requests a fresh authoritative snapshot.
export function followLogs(runId: string, reconcile: () => void, showState: (state: string) => void) {
  let socket: WebSocket | undefined;
  let timer: number | undefined;
  let failures = 0;
  let stopped = false;

  function closeSocket() {
    if (!socket) return;
    socket.onopen = socket.onmessage = socket.onclose = socket.onerror = null;
    socket.close();
    socket = undefined;
  }

  function reconnect() {
    if (stopped) return;
    window.clearTimeout(timer);
    closeSocket();
    failures++;
    showState(failures < 3
      ? 'Live connection reconnecting. Previously loaded output may be stale.'
      : 'Live connection unavailable. Retrying automatically; previously loaded output may be stale.');
    timer = window.setTimeout(connect, Math.min(1000 * 2 ** (failures - 1), 30_000));
  }

  function connect() {
    if (stopped) return;
    const url = new URL(`/api/logs/ws/${encodeURIComponent(runId)}`, location.href);
    url.protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
    try {
      socket = new WebSocket(url);
      timer = window.setTimeout(reconnect, 10_000);
      socket.onopen = () => {
        window.clearTimeout(timer);
        failures = 0;
        showState('Live connection open.');
        // Covers records persisted between the initial read and subscription.
        reconcile();
      };
      socket.onmessage = reconcile;
      socket.onclose = reconnect;
      socket.onerror = reconnect;
    } catch {
      reconnect();
    }
  }

  showState('Live connection connecting…');
  connect();
  return () => {
    stopped = true;
    window.clearTimeout(timer);
    closeSocket();
  };
}
