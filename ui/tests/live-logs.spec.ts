import { expect, type WebSocketRoute } from '@playwright/test';
import { addRun, publishLog, resetRuns, service, test, waitForCollector, waitForSavedRecords } from './fixtures';

test.beforeEach(resetRuns);

test('live logs follow broker-published records and retain identical messages once per persisted ID', async ({ page, request }) => {
  await addRun('live-output', 'Running', 1794225600);
  await publishLog('live-output', 'Saved before opening.');
  await waitForSavedRecords(request, 'live-output', 1);
  const connected = page.waitForEvent('websocket');
  await page.goto('/runs/live-output');
  const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
  await expect(logs).toContainText('Saved before opening.');
  await expect(logs.getByText('Live connection open.', { exact: true })).toBeVisible();
  const socket = await connected;
  const message = socket.waitForEvent('framereceived');
  await publishLog('live-output', 'Repeated live output', { timestamp: '2026-11-09T12:00:01Z' });
  await message;
  await publishLog('live-output', 'Repeated live output', { timestamp: '2026-11-09T12:00:01Z' });
  await expect(logs.getByRole('listitem')).toHaveText([
    '2026-11-09T12:00:00Z INFO Saved before opening.',
    '2026-11-09T12:00:01Z INFO Repeated live output',
    '2026-11-09T12:00:01Z INFO Repeated live output',
  ]);
  await logs.getByRole('button', { name: 'Refresh saved logs' }).click();
  await expect(logs.getByRole('button', { name: 'Refresh saved logs' })).toBeEnabled();
  await expect(logs.getByRole('listitem')).toHaveCount(3);
});

test('auto-scroll pauses viewport movement while ingestion and stable reading continue', async ({ page, request }, testInfo) => {
  await addRun('scroll-output', 'Running', 1794225600);
  const history = Array.from({ length: 100 }, (_, index) => `Compiler line ${index + 1}`).join('\n');
  await publishLog('scroll-output', history);
  await waitForSavedRecords(request, 'scroll-output', 1);
  await page.goto('/runs/scroll-output');
  const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
  const output = logs.getByRole('list', { name: 'Saved log records' });
  const toggle = logs.getByRole('checkbox', { name: 'Auto-scroll' });
  await expect(toggle).toBeChecked();
  await expect(logs.getByRole('listitem')).toHaveCount(1);
  const distanceFromTail = () => output.evaluate((node) => node.scrollHeight - node.clientHeight - node.scrollTop);
  await expect.poll(distanceFromTail).toBeLessThanOrEqual(1);
  await toggle.focus();
  await page.keyboard.press('Space');
  await expect(toggle).not.toBeChecked();
  await output.evaluate((node) => { node.scrollTop = 180; });
  const oldPosition = await output.evaluate((node) => node.scrollTop);
  const historyText = logs.getByText(history, { exact: true });
  await historyText.evaluate((node) => {
    const range = document.createRange();
    range.setStart(node.firstChild!, 0);
    range.setEnd(node.firstChild!, 15);
    getSelection()?.removeAllRanges();
    getSelection()?.addRange(range);
  });
  const selection = await page.evaluate(() => getSelection()?.toString());
  await publishLog('scroll-output', 'Tail while paused.', { timestamp: '2026-11-09T12:00:02Z' });
  await expect(logs.getByRole('listitem')).toHaveCount(2);
  expect(await output.evaluate((node) => node.scrollTop)).toBe(oldPosition);
  expect(await page.evaluate(() => getSelection()?.toString())).toBe(selection);
  await expect(toggle).toBeFocused();
  const historyTop = await historyText.evaluate((node) => node.getBoundingClientRect().top);
  // A late older record sorts above the reader, so retaining scrollTop alone is insufficient.
  await publishLog('scroll-output', 'Earlier output\n'.repeat(20), { timestamp: '2026-11-09T11:59:00Z' });
  await expect(logs.getByRole('listitem')).toHaveCount(3);
  expect(Math.abs(await historyText.evaluate((node) => node.getBoundingClientRect().top) - historyTop)).toBeLessThanOrEqual(1);
  expect(await page.evaluate(() => getSelection()?.toString())).toBe(selection);
  await page.keyboard.press('Space');
  await expect(toggle).toBeChecked();
  await expect.poll(distanceFromTail).toBeLessThanOrEqual(1);
  await publishLog('scroll-output', 'Current tail\n'.repeat(20), { timestamp: '2026-11-09T12:00:03Z' });
  await expect(logs.getByRole('listitem')).toHaveCount(4);
  await expect.poll(distanceFromTail).toBeLessThanOrEqual(1);
  await expect(toggle).toBeFocused();
  for (const width of [1280, 390]) {
    await page.setViewportSize({ width, height: 900 });
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
    await page.screenshot({ path: testInfo.outputPath(`live-logs-${width}.png`), fullPage: true });
  }
});

test('records persisted during initial history loading are reconciled after the delayed snapshot', async ({ page, request }) => {
  await addRun('overlap-output', 'Running', 1794225600);
  await publishLog('overlap-output', 'Before the snapshot.');
  await waitForSavedRecords(request, 'overlap-output', 1);
  let release!: () => void;
  let received!: () => void;
  const gate = new Promise<void>((resolve) => { release = resolve; });
  const waiting = new Promise<void>((resolve) => { received = resolve; });
  let delayed = false;
  await page.route('**/api/logs/overlap-output', async (route) => {
    if (delayed) { await route.continue(); return; }
    delayed = true;
    const response = await route.fetch();
    received();
    await gate;
    await route.fulfill({ response });
  });
  try {
    const connected = page.waitForEvent('websocket');
    await page.goto('/runs/overlap-output');
    await waiting;
    const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
    await expect(logs).toContainText('Loading saved logs…');
    await expect(logs).toContainText('Live connection open.');
    const socket = await connected;
    const frame = socket.waitForEvent('framereceived');
    await publishLog('overlap-output', 'During the snapshot.', { timestamp: '2026-11-09T12:00:01Z' });
    await frame;
    await waitForSavedRecords(request, 'overlap-output', 2);
    release();
    await expect(logs.getByRole('listitem')).toHaveText([
      '2026-11-09T12:00:00Z INFO Before the snapshot.',
      '2026-11-09T12:00:01Z INFO During the snapshot.',
    ]);
  } finally {
    release();
    await page.unrouteAll({ behavior: 'wait' });
  }
});

test('interrupted live transport retains output and recovers persisted repeated records in order', async ({ page, request }) => {
  await addRun('reconnect-output', 'Running', 1794225600);
  await publishLog('reconnect-output', 'Retained history.');
  await waitForSavedRecords(request, 'reconnect-output', 1);
  await page.clock.install();
  let interrupted = false;
  let connection: WebSocketRoute | undefined;
  let upstream: WebSocketRoute | undefined;
  // Interrupt transport only. All delivered frames and snapshots still come from real services.
  await page.routeWebSocket('**/api/logs/ws/reconnect-output', (route) => {
    if (interrupted) { void route.close(); return; }
    connection = route;
    upstream = route.connectToServer();
  });
  await page.goto('/runs/reconnect-output');
  const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
  await expect(logs).toContainText('Live connection open.');
  await expect(logs.getByRole('listitem')).toHaveCount(1);
  interrupted = true;
  await connection!.close();
  await upstream!.close();
  await expect(logs).toContainText('Live connection reconnecting.');
  await expect(logs).toContainText('Retained history.');
  await page.clock.runFor(1000);
  await expect(logs).toContainText('Live connection reconnecting.');
  await page.clock.runFor(2000);
  await expect(logs).toContainText('Live connection unavailable.');
  await publishLog('reconnect-output', 'Repeated during disconnect.', { timestamp: '2026-11-09T12:00:02Z' });
  await publishLog('reconnect-output', 'Repeated during disconnect.', { timestamp: '2026-11-09T12:00:02Z' });
  await publishLog('reconnect-output', 'Earlier persisted output.', { timestamp: '2026-11-09T12:00:01Z' });
  await waitForSavedRecords(request, 'reconnect-output', 4);
  interrupted = false;
  await page.clock.runFor(4000);
  await expect(logs).toContainText('Live connection open.');
  await expect(logs.getByRole('listitem')).toHaveText([
    '2026-11-09T12:00:00Z INFO Retained history.',
    '2026-11-09T12:00:01Z INFO Earlier persisted output.',
    '2026-11-09T12:00:02Z INFO Repeated during disconnect.',
    '2026-11-09T12:00:02Z INFO Repeated during disconnect.',
  ]);
  await publishLog('reconnect-output', 'Live again.', { timestamp: '2026-11-09T12:00:03Z' });
  await expect(logs.getByRole('listitem')).toHaveCount(5);
});

for (const status of ['Completed', 'Failed', 'Canceled']) {
  test(`late persisted logs remain discoverable after a run becomes ${status}`, async ({ page, request }) => {
    const id = `terminal-${status}`;
    await addRun(id, 'Running', 1794225600);
    await publishLog(id, 'Before completion.');
    await waitForSavedRecords(request, id, 1);
    await page.clock.install();
    let forward = true;
    await page.routeWebSocket(`**/api/logs/ws/${id}`, (route) => {
      const server = route.connectToServer();
      server.onMessage((message) => { if (forward) route.send(message); });
    });
    await page.goto(`/runs/${id}`);
    const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
    const details = page.getByRole('region', { name: 'Run details', exact: true });
    const toggle = logs.getByRole('checkbox', { name: 'Auto-scroll' });
    await expect(logs.getByRole('listitem')).toHaveCount(1);
    await expect(logs).toContainText('Live connection open.');
    await toggle.uncheck();
    await toggle.focus();
    const response = await request.put(`http://127.0.0.1:${process.env.API_TEST_PORT}/app/cicd/api/runs/${id}`, {
      data: { status, time_started: 1794225610, time_ended: 1794225675 },
    });
    expect(response.ok()).toBe(true);
    await publishLog(id, 'Around completion.', { timestamp: '2026-11-09T12:00:01Z' });
    await page.clock.runFor(5000);
    await expect(details.getByText(status, { exact: true })).toBeVisible();
    await expect(logs.getByRole('listitem')).toHaveCount(2);
    // Drop live notifications to prove periodic reads continue after terminal status.
    forward = false;
    await publishLog(id, 'Asynchronously ingested final output.', { timestamp: '2026-11-09T12:00:02Z' });
    await waitForSavedRecords(request, id, 3);
    await page.clock.runFor(5000);
    await expect(logs.getByRole('listitem')).toHaveCount(3);
    await page.clock.runFor(60_000);
    await publishLog(id, 'Much later output.', { timestamp: '2026-11-09T12:00:03Z' });
    await waitForSavedRecords(request, id, 4);
    await page.clock.runFor(5000);
    await expect(logs.getByRole('listitem')).toHaveCount(4);
    await expect(logs.getByRole('listitem').last()).toContainText('Much later output.');
    await expect(toggle).not.toBeChecked();
    await expect(toggle).toBeFocused();
  });
}

test('switching runs closes the old live subscription and cancels its scheduled reads', async ({ page, request }) => {
  for (const id of ['live-run-a', 'live-run-b']) {
    await addRun(id, 'Running', 1794225600);
    await publishLog(id, `Output for ${id}.`);
    await waitForSavedRecords(request, id, 1);
  }
  await page.clock.install();
  const subscriptions: string[] = [];
  const reads: string[] = [];
  page.on('websocket', (socket) => subscriptions.push(socket.url()));
  page.on('request', (request) => { if (request.url().includes('/api/logs/')) reads.push(request.url()); });
  await page.goto('/');
  const connected = page.waitForEvent('websocket');
  await page.getByRole('link', { name: 'Open run live-run-a', exact: true }).click();
  const socket = await connected;
  const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
  await expect(logs).toContainText('Output for live-run-a.');
  const closed = socket.waitForEvent('close');
  await page.getByRole('link', { name: 'Back to runs', exact: true }).click();
  await closed;
  await page.getByRole('link', { name: 'Open run live-run-b', exact: true }).click();
  await expect(logs).toContainText('Output for live-run-b.');
  const oldReads = reads.filter((url) => url.endsWith('/live-run-a')).length;
  await publishLog('live-run-a', 'Old run output after switching.', { timestamp: '2026-11-09T12:00:01Z' });
  await page.clock.runFor(31_000);
  await expect(logs.getByRole('listitem')).toHaveCount(1);
  await expect(logs).not.toContainText('live-run-a');
  expect(reads.filter((url) => url.endsWith('/live-run-a'))).toHaveLength(oldReads);
  expect(subscriptions.filter((url) => url.endsWith('/live-run-a'))).toHaveLength(1);
});

test('collector outage retries automatically while readable details and log controls remain usable', async ({ page, request }) => {
  await addRun('collector-recovery', 'Running', 1794225600);
  await publishLog('collector-recovery', 'Retained during outage.');
  await waitForSavedRecords(request, 'collector-recovery', 1);
  await page.clock.install();
  await page.goto('/runs/collector-recovery');
  const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
  await expect(logs.getByRole('listitem')).toHaveCount(1);
  await expect(logs).toContainText('Live connection open.');
  const toggle = logs.getByRole('checkbox', { name: 'Auto-scroll' });
  await toggle.uncheck();
  await toggle.focus();
  service('stop', 'logcollector');
  try {
    await expect(logs).toContainText('Live connection reconnecting.');
    await page.clock.runFor(1000);
    await page.clock.runFor(2000);
    await expect(logs).toContainText('Live connection unavailable.');
    await page.clock.runFor(5000);
    await expect(logs.getByRole('alert')).toContainText('stale');
    await expect(logs).toContainText('Retained during outage.');
    await expect(page.getByRole('region', { name: 'Run details', exact: true }).getByText('Running', { exact: true })).toBeVisible();
    await expect(toggle).not.toBeChecked();
    await expect(toggle).toBeFocused();
    await publishLog('collector-recovery', 'Queued while collector was down.', { timestamp: '2026-11-09T12:00:01Z' });
    service('start', 'logcollector');
    await waitForCollector(request);
    await waitForSavedRecords(request, 'collector-recovery', 2);
    // Allow real handshake I/O between timer steps instead of expiring it in one clock jump.
    await expect.poll(async () => {
      await page.clock.runFor(1000);
      return logs.getByText('Live connection open.', { exact: true }).isVisible();
    }, { intervals: [100], timeout: 10_000 }).toBe(true);
    await expect(logs.getByRole('listitem')).toHaveCount(2);
    await expect(logs.getByRole('alert')).toHaveCount(0);
    await expect(toggle).not.toBeChecked();
    await expect(toggle).toBeFocused();
  } finally {
    service('start', 'logcollector');
  }
});

test('run API failure is independent of continuing live logs and preserves viewer controls', async ({ page, request }) => {
  await addRun('independent-output', 'Running', 1794225600);
  await publishLog('independent-output', 'Before the API outage.');
  await waitForSavedRecords(request, 'independent-output', 1);
  await page.clock.install();
  await page.goto('/runs/independent-output');
  const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
  await expect(logs.getByRole('listitem')).toHaveCount(1);
  await expect(logs).toContainText('Live connection open.');
  const toggle = logs.getByRole('checkbox', { name: 'Auto-scroll' });
  await toggle.uncheck();
  await toggle.focus();
  service('stop', 'api');
  try {
    await page.clock.runFor(5000);
    await expect(page.locator('#notice[role="alert"]')).toContainText('stale');
    await publishLog('independent-output', 'Output continues during the API outage.', { timestamp: '2026-11-09T12:00:01Z' });
    await expect(logs.getByRole('listitem')).toHaveCount(2);
    await expect(logs.getByRole('alert')).toHaveCount(0);
    await expect(logs).toContainText('Live connection open.');
    await expect(toggle).not.toBeChecked();
    await expect(toggle).toBeFocused();
    service('start', 'api');
    await expect.poll(async () => (await request.get('/api/runs')).status()).toBe(200);
    await page.clock.runFor(5000);
    await expect(page.locator('#notice[role="alert"]')).toHaveCount(0);
    await expect(logs.getByRole('listitem')).toHaveCount(2);
    await expect(toggle).not.toBeChecked();
    await expect(toggle).toBeFocused();
  } finally {
    service('start', 'api');
  }
});
