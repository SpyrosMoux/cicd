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
        <p class="description">Build history across your repositories.</p></div>
      <button id="refresh" type="button">Refresh runs <span aria-hidden="true">↻</span></button>
    </div>
    <div class="list-heading"><h2>All repositories</h2><span id="run-count"></span></div>
    <div id="notice" role="status" aria-live="polite"></div>
    <div id="content" class="panel"></div>
    <footer>FlowForge <span aria-hidden="true">/</span> Your pipelines, in view.</footer>
  </section>`;

const content = document.querySelector<HTMLElement>('#content')!;
const notice = document.querySelector<HTMLElement>('#notice')!;
const count = document.querySelector<HTMLElement>('#run-count')!;
const refresh = document.querySelector<HTMLButtonElement>('#refresh')!;
const statuses = new Set(['Pending', 'Running', 'Completed', 'Failed', 'Canceled']);
const durations = new Map<HTMLElement, Run>();
let hasSnapshot = false;

function renderRun(run: Run) {
  const row = element('tr');
  const identity = element('td', 'identity');
  identity.append(element('strong', 'repository', run.repository || 'Unavailable'));
  identity.append(element('span', 'branch', run.branch || 'Unavailable'));
  identity.append(element('code', 'run-id', run.id));
  const status = element('td');
  let statusClass = 'unknown';
  let statusLabel = 'Unknown';
  if (statuses.has(run.status)) {
    statusClass = run.status.toLowerCase();
    statusLabel = run.status;
  } else if (run.status) {
    statusLabel = `Unknown (${run.status})`;
  }
  status.append(element('span', `status ${statusClass}`, statusLabel));
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

function renderRuns(runs: Run[]) {
  durations.clear();
  if (runs.length === 0) {
    const state = element('div', 'state');
    const icon = element('div', 'state-icon', '◇');
    icon.setAttribute('aria-hidden', 'true');
    state.append(icon, element('h2', '', 'No runs yet'),
      element('p', '', 'Pipeline runs will appear here when your repositories send work to FlowForge.'));
    content.replaceChildren(state);
    return;
  }
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
  for (const run of runs) {
    body.append(renderRun(run));
  }
  table.append(head, body);
  content.replaceChildren(table);
}

async function refreshRuns() {
  if (refresh.disabled) return;
  refresh.disabled = true;
  notice.setAttribute('role', 'status');
  notice.textContent = 'Loading runs…';
  content.setAttribute('aria-busy', 'true');
  try {
    const runs = await readRuns();
    renderRuns(runs);
    hasSnapshot = true;
    count.textContent = `${runs.length} ${runs.length === 1 ? 'run' : 'runs'}`;
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
  }
}

refresh.addEventListener('click', refreshRuns);
window.setInterval(() => {
  for (const [node, run] of durations) node.textContent = `Duration: ${formatDuration(run)}`;
}, 1000);
await refreshRuns();
