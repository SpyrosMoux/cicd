import { element } from './dom';
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
  private readonly output = element('ol', 'log-output');
  private readonly refresh = element('button', '', 'Refresh saved logs');
  private run: Run | null = null;
  private request?: AbortController;
  private records: LogRecord[] = [];
  private hasSnapshot = false;
  private failed = false;

  constructor(private readonly host: HTMLElement) {
    const heading = element('h2', '', 'Saved logs');
    heading.id = 'logs-heading';
    const toolbar = element('div', 'log-toolbar');
    this.refresh.type = 'button';
    this.refresh.addEventListener('click', () => this.load());
    toolbar.append(heading, this.refresh);
    this.notice.setAttribute('role', 'status');
    this.notice.setAttribute('aria-live', 'polite');
    this.output.setAttribute('aria-label', 'Saved log records');
    this.output.tabIndex = 0;
    this.host.append(toolbar, this.notice, this.output);
  }

  select(run: Run | null) {
    if (run?.id === this.run?.id) {
      this.run = run;
      if (this.hasSnapshot && !this.failed && !this.request) this.showSnapshotNotice();
      return;
    }
    this.request?.abort();
    this.request = undefined;
    this.run = run;
    this.records = [];
    this.hasSnapshot = false;
    this.failed = false;
    this.output.replaceChildren();
    this.host.hidden = run === null;
    if (run) void this.load();
  }

  private showSnapshotNotice() {
    this.notice.setAttribute('role', 'status');
    if (this.records.length) {
      this.notice.textContent = `${this.records.length} saved ${this.records.length === 1 ? 'record' : 'records'}.`;
    } else {
      this.notice.textContent = this.run && ['Pending', 'Running'].includes(this.run.status)
        ? 'No saved logs yet. Waiting for output; refresh to check again.'
        : 'No saved logs for this run.';
    }
  }

  private async load() {
    if (!this.run || this.request) return;
    const request = new AbortController();
    this.request = request;
    this.refresh.disabled = true;
    this.output.setAttribute('aria-busy', 'true');
    this.notice.setAttribute('role', 'status');
    this.notice.textContent = this.hasSnapshot
      ? 'Refreshing saved logs. Previously loaded output may be stale.' : 'Loading saved logs…';
    try {
      const records = await readLogs(this.run.id, request.signal);
      if (this.request !== request) return;
      this.records = records;
      this.hasSnapshot = true;
      this.failed = false;
      this.output.replaceChildren(...records.map((record) => {
        const line = element('li', 'log-record');
        line.dataset.logId = record.id;
        const timestamp = element('span', 'log-timestamp', record.timestamp || 'Timestamp unavailable');
        const severity = ['INFO', 'WARN', 'ERROR', 'DEBUG'][record.logLevel] ?? `Level ${record.logLevel}`;
        line.append(timestamp, ' ', element('span', 'log-level', severity), ' ', element('span', 'log-message', record.message));
        return line;
      }));
      this.showSnapshotNotice();
    } catch {
      if (this.request !== request) return;
      this.failed = true;
      this.notice.setAttribute('role', 'alert');
      this.notice.textContent = this.hasSnapshot
        ? 'Unable to refresh saved logs. Previously loaded output is stale. Use Refresh saved logs to try again.'
        : 'The log collector is unavailable. Saved logs could not be retrieved. Use Refresh saved logs to try again.';
    } finally {
      if (this.request === request) {
        this.request = undefined;
        this.refresh.disabled = false;
        this.output.setAttribute('aria-busy', 'false');
      }
    }
  }
}
