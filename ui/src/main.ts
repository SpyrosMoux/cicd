import { readRuns, formatTime, formatDuration, type Run } from './runs';
import './style.css';

function element<K extends keyof HTMLElementTagNameMap>(tag: K, className = '', text = '') {
  const node = document.createElement(tag);
  node.className = className;
  node.textContent = text;
  return node;
}

const app = document.querySelector<HTMLElement>('#app')!;
app.innerHTML = `
  <header class="topbar">
    <a class="brand" href="/" aria-label="FlowForge home"><span class="brand-mark" aria-hidden="true">F</span>FlowForge</a>
    <span class="mode">Read-only dashboard</span>
  </header>
  <section class="workspace" aria-labelledby="page-title">
    <div class="intro">
      <div><p class="eyebrow">PIPELINE ACTIVITY</p><h1 id="page-title">Recent runs</h1>
        <p class="description">Build history across your repositories. Updates automatically.</p></div>
      <button id="refresh" type="button">Refresh runs <span aria-hidden="true">↻</span></button>
    </div>
    <div class="filters" role="group" aria-label="Run filters">
      <label for="repository-filter">Repository<select id="repository-filter"><option value="">All repositories</option></select></label>
      <label for="status-filter">Status<select id="status-filter"><option value="">All statuses</option></select></label>
      <button id="clear-filters" type="button" disabled>Clear filters</button>
    </div>
    <div class="list-heading"><h2>Run history</h2><span id="run-count"></span></div>
    <div id="notice" role="status" aria-live="polite"></div>
    <div id="content" class="panel"></div>
    <footer>FlowForge <span aria-hidden="true">/</span> Your pipelines, in view.</footer>
  </section>`;

const content = document.querySelector<HTMLElement>('#content')!;
const notice = document.querySelector<HTMLElement>('#notice')!;
const count = document.querySelector<HTMLElement>('#run-count')!;
const refresh = document.querySelector<HTMLButtonElement>('#refresh')!;
const repositoryFilter = document.querySelector<HTMLSelectElement>('#repository-filter')!;
const statusFilter = document.querySelector<HTMLSelectElement>('#status-filter')!;
const clearFilters = document.querySelector<HTMLButtonElement>('#clear-filters')!;
const statuses = new Set(['Pending', 'Running', 'Completed', 'Failed', 'Canceled']);
const durations = new Map<HTMLElement, Run>();
const renderedRows = new Map<string, { signature: string; node: HTMLTableRowElement }>();
let hasSnapshot = false;
let snapshot: Run[] = [];
let refreshTimer: number | undefined;

function updateOptions(select: HTMLSelectElement, values: string[], label: (value: string) => string) {
  const selected = select.value;
  // Keep an active selection available even when its last run leaves the snapshot.
  if (selected && !values.includes(selected)) values.push(selected);
  const unchanged = select.options.length === values.length + 1 &&
    values.every((value, index) => select.options[index + 1].value === value);
  if (!unchanged) {
    select.replaceChildren(select.options[0], ...values.map((value) => new Option(label(value), value)));
    select.value = selected;
  }
}

function statusLabel(status: string) {
  if (statuses.has(status)) return status;
  return status ? `Unknown (${status})` : 'Unknown';
}

function renderSnapshot() {
  const filtered = repositoryFilter.value !== '' || statusFilter.value !== '';
  const runs = snapshot.filter((run) =>
    (!repositoryFilter.value || run.repository === repositoryFilter.value) &&
    (!statusFilter.value || run.status === statusFilter.value));
  clearFilters.disabled = !filtered;
  if (filtered) {
    count.textContent = `${runs.length} of ${snapshot.length} runs`;
  } else {
    count.textContent = `${runs.length} ${runs.length === 1 ? 'run' : 'runs'}`;
  }
  renderRuns(runs);
}

function renderRun(run: Run) {
  const row = element('tr');
  const identity = element('td', 'identity');
  identity.append(element('strong', 'repository', run.repository || 'Unavailable'));
  identity.append(element('span', 'branch', run.branch || 'Unavailable'));
  identity.append(element('code', 'run-id', run.id));
  const status = element('td');
  const statusClass = statuses.has(run.status) ? run.status.toLowerCase() : 'unknown';
  status.append(element('span', `status ${statusClass}`, statusLabel(run.status)));
  const trigger = element('td', 'trigger');
  trigger.append(element('span', '', run.trigger_type || 'Unavailable'));
  trigger.append(element('span', 'muted', run.triggered_by || 'Unavailable'));
  const timing = element('td', 'timing');
  timing.append(element('span', '', `Triggered: ${formatTime(run.time_triggered)}`));
  let start = formatTime(run.time_started);
  let end = formatTime(run.time_ended);
  if (start === 'Unavailable' && run.status === 'Pending') start = 'Pending';
  if (end === 'Unavailable' && run.status === 'Running') end = 'In progress';
  timing.append(element('span', '', `Started: ${start}`));
  timing.append(element('span', '', `Ended: ${end}`));
  const duration = element('span', '', `Duration: ${formatDuration(run)}`);
  timing.append(duration);
  if (run.status === 'Running') durations.set(duration, run);
  row.append(identity, status, trigger, timing);
  return row;
}

function createRunTable() {
  const table = element('table');
  table.append(element('caption', 'sr-only', 'Recent pipeline runs, newest triggered first'));
  const head = element('thead');
  const headings = element('tr');
  for (const label of ['Repository / run', 'Status', 'Trigger', 'Timing']) {
    const cell = element('th', '', label);
    cell.scope = 'col';
    headings.append(cell);
  }
  head.append(headings);
  const body = element('tbody');
  table.append(head, body);
  content.replaceChildren(table);
  return body;
}

function renderEmptyRuns() {
  durations.clear();
  renderedRows.clear();
  const state = element('div', 'state');
  const icon = element('div', 'state-icon', '◇');
  icon.setAttribute('aria-hidden', 'true');
  const noMatches = snapshot.length > 0;
  state.append(icon, element('h2', '', noMatches ? 'No matching runs' : 'No runs yet'),
    element('p', '', noMatches ? 'Choose different filters or clear them to see all runs.'
      : 'Pipeline runs will appear here when your repositories send work to FlowForge.'));
  content.replaceChildren(state);
}

function renderRuns(runs: Run[]) {
  if (runs.length === 0) {
    renderEmptyRuns();
    return;
  }
  const body = content.querySelector('tbody') ?? createRunTable();
  const nextRows = runs.map((run) => {
    const signature = JSON.stringify(run);
    let row = renderedRows.get(run.id);
    if (row?.signature !== signature) {
      row = { signature, node: renderRun(run) };
      renderedRows.set(run.id, row);
    }
    return row.node;
  });
  const nextNodes = new Set(nextRows);
  // Anchor an unchanged visible row when new work is inserted above the viewport.
  const anchor = window.scrollY > 0 ? Array.from(body.rows).find((row) => {
    const bounds = row.getBoundingClientRect();
    return nextNodes.has(row) && bounds.bottom > 0 && bounds.top < window.innerHeight;
  }) : undefined;
  const anchorTop = anchor?.getBoundingClientRect().top ?? 0;
  for (const row of Array.from(body.rows)) {
    if (!nextNodes.has(row)) row.remove();
  }
  nextRows.forEach((row, index) => {
    if (body.rows[index] !== row) body.insertBefore(row, body.rows[index] ?? null);
  });
  if (anchor) window.scrollBy(0, anchor.getBoundingClientRect().top - anchorTop);
  const visibleIds = new Set(runs.map((run) => run.id));
  for (const id of renderedRows.keys()) {
    if (!visibleIds.has(id)) renderedRows.delete(id);
  }
  for (const node of durations.keys()) {
    if (!node.isConnected) durations.delete(node);
  }
}

async function refreshRuns() {
  if (refresh.disabled) return;
  window.clearTimeout(refreshTimer);
  refresh.disabled = true;
  if (!hasSnapshot) {
    notice.setAttribute('role', 'status');
    notice.textContent = 'Loading runs…';
  }
  content.setAttribute('aria-busy', 'true');
  try {
    const runs = await readRuns();
    snapshot = runs;
    hasSnapshot = true;
    updateOptions(repositoryFilter, [...new Set(runs.map((run) => run.repository).filter(Boolean))]
      .sort((a, b) => a.localeCompare(b)), (value) => value);
    const extraStatuses = [...new Set(runs.map((run) => run.status)
      .filter((value) => value && !statuses.has(value)))].sort((a, b) => a.localeCompare(b));
    updateOptions(statusFilter, [...statuses, ...extraStatuses], statusLabel);
    renderSnapshot();
    notice.setAttribute('role', 'status');
    notice.textContent = `Last refreshed: ${formatTime(Date.now() / 1000)}`;
  } catch {
    notice.setAttribute('role', 'alert');
    notice.textContent = hasSnapshot
      ? 'Unable to refresh runs. Showing previously loaded history; it may be stale.'
      : 'The run API is unavailable. History could not be retrieved.';
    const retry = element('button', '', 'Try again');
    retry.type = 'button';
    retry.addEventListener('click', refreshRuns);
    if (hasSnapshot) {
      notice.append(' ', retry);
    } else {
      const state = element('div', 'state');
      state.append(element('h2', '', 'Unable to load runs'),
        element('p', '', 'Check that the API is running and reachable, then try again.'), retry);
      content.replaceChildren(state);
      count.textContent = 'Unavailable';
    }
  } finally {
    refresh.disabled = false;
    content.setAttribute('aria-busy', 'false');
    // Schedule after completion so slow requests and manual refreshes never overlap.
    refreshTimer = window.setTimeout(refreshRuns, 5000);
  }
}

repositoryFilter.addEventListener('change', () => { if (hasSnapshot) renderSnapshot(); });
statusFilter.addEventListener('change', () => { if (hasSnapshot) renderSnapshot(); });
clearFilters.addEventListener('click', () => {
  repositoryFilter.value = '';
  statusFilter.value = '';
  renderSnapshot();
  repositoryFilter.focus();
});
refresh.addEventListener('click', refreshRuns);
window.setInterval(() => {
  for (const [node, run] of durations) node.textContent = `Duration: ${formatDuration(run)}`;
}, 1000);
await refreshRuns();
