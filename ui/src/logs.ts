import { element } from './dom';
import { followLogs } from './live';
import type { Run } from './runs';

interface LogRecord {
  id: string;
  runid: string;
  timestamp: string;
  logLevel: number;
  message: string;
}

function isLogRecord(value: unknown, runId: string): value is LogRecord {
  if (typeof value !== 'object' || value === null) return false;
  const record = value as Record<string, unknown>;
  return typeof record.id === 'string' && record.id !== '' && record.runid === runId &&
    typeof record.timestamp === 'string' && typeof record.message === 'string' &&
    typeof record.logLevel === 'number' && Number.isInteger(record.logLevel);
}

async function readLogs(runId: string, signal: AbortSignal): Promise<LogRecord[]> {
  const response = await fetch(`/api/logs/${encodeURIComponent(runId)}`, {
    cache: 'no-store', signal: AbortSignal.any([signal, AbortSignal.timeout(10_000)]),
  });
  if (!response.ok) throw new Error('The log collector is unavailable.');
  const envelope = await response.json();
  if (envelope?.status !== 200 || envelope.error !== '' ||
      !(envelope.data === null || (Array.isArray(envelope.data) &&
        envelope.data.every((value: unknown) => isLogRecord(value, runId))))) {
    throw new Error('The log collector returned an unexpected response.');
  }
  // The existing collector does not order reads. Sort the snapshot without inventing emission order.
  const records: LogRecord[] = envelope.data ?? [];
  return records.sort((a, b) => {
    if (a.timestamp !== b.timestamp) return a.timestamp < b.timestamp ? -1 : 1;
    if (a.id === b.id) return 0;
    return a.id < b.id ? -1 : 1;
  });
}

export class SavedLogs {
  private readonly notice = element('p', 'log-notice');
  private readonly connection = element('p', 'log-notice');
  private readonly output = element('ol', 'log-output');
  private readonly refresh = element('button', '', 'Refresh saved logs');
  private readonly autoScroll = element('input');
  private run: Run | null = null;
  private request?: AbortController;
  private records: LogRecord[] = [];
  private hasSnapshot = false;
  private failed = false;
  private timer?: number;
  private stopFollowing?: () => void;
  private reconcilePending = false;
  private failures = 0;
  private readonly lines = new Map<string, { record: string; node: HTMLLIElement }>();

  constructor(private readonly host: HTMLElement) {
    const heading = element('h2', '', 'Saved logs');
    heading.id = 'logs-heading';
    const toolbar = element('div', 'log-toolbar');
    this.refresh.type = 'button';
    this.refresh.addEventListener('click', () => this.load());
    const scrollControl = element('label', 'log-autoscroll');
    this.autoScroll.type = 'checkbox';
    this.autoScroll.checked = true;
    this.autoScroll.addEventListener('change', () => {
      if (this.autoScroll.checked) this.followTail();
    });
    scrollControl.append(this.autoScroll, 'Auto-scroll');
    toolbar.append(heading, scrollControl, this.refresh);
    this.notice.setAttribute('role', 'status');
    this.notice.setAttribute('aria-live', 'polite');
    this.output.setAttribute('aria-label', 'Saved log records');
    this.output.tabIndex = 0;
    this.connection.setAttribute('aria-live', 'polite');
    this.host.append(toolbar, this.connection, this.notice, this.output);
  }

  select(run: Run | null) {
    if (run?.id === this.run?.id) {
      const becameTerminal = run && run.status !== this.run?.status &&
        ['Completed', 'Failed', 'Canceled'].includes(run.status);
      this.run = run;
      if (becameTerminal) void this.load();
      if (this.hasSnapshot && !this.failed && !this.request) this.showSnapshotNotice();
      return;
    }
    this.stopFollowing?.();
    window.clearTimeout(this.timer);
    this.request?.abort();
    this.request = undefined;
    this.run = run;
    this.records = [];
    this.lines.clear();
    this.reconcilePending = false;
    this.failures = 0;
    this.hasSnapshot = false;
    this.failed = false;
    this.output.replaceChildren();
    this.host.hidden = run === null;
    if (run) {
      this.stopFollowing = followLogs(run.id, () => { void this.load(); }, (state) => {
        this.connection.textContent = state;
      });
      void this.load();
    }
  }

  private showSnapshotNotice() {
    this.notice.setAttribute('role', 'status');
    if (this.records.length) {
      this.notice.textContent = `${this.records.length} saved ${this.records.length === 1 ? 'record' : 'records'}.`;
    } else {
      this.notice.textContent = this.run && ['Pending', 'Running'].includes(this.run.status)
        ? 'No saved logs yet. Waiting for output.'
        : 'No saved logs for this run.';
    }
  }

  private retryDelay() {
    return this.failed ? Math.min(1000 * 2 ** (this.failures - 1), 30_000) : 5000;
  }

  private followTail() {
    this.output.scrollTop = this.output.scrollHeight;
  }

  private renderRecords() {
    const scroll = this.output.scrollTop;
    const bounds = this.output.getBoundingClientRect();
    const anchor = Array.from(this.output.children).find((line) => {
      const lineBounds = line.getBoundingClientRect();
      return lineBounds.bottom > bounds.top && lineBounds.top < bounds.bottom;
    });
    const anchorTop = anchor?.getBoundingClientRect().top ?? 0;
    const nextIds = new Set(this.records.map((record) => record.id));
    for (const [id, line] of this.lines) {
      if (!nextIds.has(id)) {
        line.node.remove();
        this.lines.delete(id);
      }
    }
    this.records.forEach((record, index) => {
      const signature = JSON.stringify(record);
      let line = this.lines.get(record.id);
      if (line?.record !== signature) {
        const node = element('li', 'log-record');
        node.dataset.logId = record.id;
        const timestamp = element('span', 'log-timestamp', record.timestamp || 'Timestamp unavailable');
        const severity = ['INFO', 'WARN', 'ERROR', 'DEBUG'][record.logLevel] ?? `Level ${record.logLevel}`;
        node.append(timestamp, ' ', element('span', 'log-level', severity), ' ', element('span', 'log-message', record.message));
        line?.node.remove();
        line = { record: signature, node };
        this.lines.set(record.id, line);
      }
      if (this.output.children[index] !== line.node) {
        this.output.insertBefore(line.node, this.output.children[index] ?? null);
      }
    });
    if (this.autoScroll.checked) {
      this.followTail();
    } else if (anchor?.isConnected) {
      this.output.scrollTop = scroll + anchor.getBoundingClientRect().top - anchorTop;
    } else {
      this.output.scrollTop = scroll;
    }
  }

  private async load() {
    if (!this.run) return;
    if (this.request) {
      this.reconcilePending = true;
      return;
    }
    window.clearTimeout(this.timer);
    const request = new AbortController();
    this.request = request;
    this.output.setAttribute('aria-busy', 'true');
    if (!this.hasSnapshot) {
      this.notice.setAttribute('role', 'status');
      this.notice.textContent = 'Loading saved logs…';
    }
    try {
      const records = await readLogs(this.run.id, request.signal);
      if (this.request !== request) return;
      this.records = [...new Map(records.map((record) => [record.id, record])).values()];
      this.hasSnapshot = true;
      this.failed = false;
      this.failures = 0;
      this.renderRecords();
      this.showSnapshotNotice();
    } catch {
      if (this.request !== request) return;
      this.failed = true;
      this.failures++;
      this.notice.setAttribute('role', 'alert');
      this.notice.textContent = this.hasSnapshot
        ? 'Unable to refresh saved logs. Previously loaded output is stale. Use Refresh saved logs to try again.'
        : 'The log collector is unavailable. Saved logs could not be retrieved. Use Refresh saved logs to try again.';
    } finally {
      if (this.request === request) {
        this.request = undefined;
        this.output.setAttribute('aria-busy', 'false');
        // Signals during a read require another snapshot, even when that read succeeds.
        const delay = this.reconcilePending ? 100 : this.retryDelay();
        this.reconcilePending = false;
        this.timer = window.setTimeout(() => { void this.load(); }, delay);
      }
    }
  }
}
