export interface Run {
  id: string;
  repository: string;
  branch: string;
  status: string;
  trigger_type: string;
  triggered_by: string;
  time_triggered?: number | null;
  time_started?: number | null;
  time_ended?: number | null;
}

function isRun(value: unknown): value is Run {
  if (typeof value !== 'object' || value === null) return false;
  const run = value as Record<string, unknown>;
  return ['id', 'repository', 'branch', 'status', 'trigger_type', 'triggered_by']
    .every((key) => typeof run[key] === 'string') &&
    ['time_triggered', 'time_started', 'time_ended']
      .every((key) => run[key] == null || (typeof run[key] === 'number' && Number.isFinite(run[key])));
}

export async function readRuns(): Promise<Run[]> {
  const response = await fetch('/api/runs', {
    cache: 'no-store', signal: AbortSignal.timeout(10_000),
  });
  if (!response.ok) throw new Error('The run API is unavailable.');
  const envelope = await response.json();
  if (envelope?.status !== 200 || envelope.error !== '' ||
      !(envelope.data === null || (Array.isArray(envelope.data) && envelope.data.every(isRun)))) {
    throw new Error('The run API returned an unexpected response.');
  }
  const runs: Run[] = envelope.data ?? [];
  return runs.sort((a, b) => {
    const triggered = (b.time_triggered || 0) - (a.time_triggered || 0);
    if (triggered !== 0) return triggered;
    if (a.id === b.id) return 0;
    return a.id < b.id ? -1 : 1;
  });
}

export function formatTime(value?: number | null): string {
  if (value == null || value <= 0) return 'Unavailable';
  const date = new Date(value * 1000);
  if (!Number.isFinite(date.getTime())) return 'Unavailable';
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: 'medium', timeStyle: 'medium', timeZone: 'UTC',
  }).format(date) + ' UTC';
}

export function formatDuration(run: Run): string {
  const start = run.time_started ?? 0;
  const end = run.status === 'Running' ? Date.now() / 1000 : run.time_ended ?? 0;
  if (formatTime(start) === 'Unavailable' || formatTime(end) === 'Unavailable' || end < start) return 'Unavailable';
  const seconds = Math.floor(end - start);
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor(seconds % 3600 / 60);
  const parts = [];
  if (hours) parts.push(`${hours}h`);
  if (minutes) parts.push(`${minutes}m`);
  parts.push(`${seconds % 60}s`);
  return parts.join(' ');
}
