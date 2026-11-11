import { expect } from '@playwright/test';
import { addRun, database, delaySavedLogResponse, publishLog, resetRuns, service, test, waitForCollector, waitForSavedRecords } from './fixtures';

test.beforeEach(resetRuns);

test('saved logs preserve repeated records, metadata and stable presentation of unordered reads', async ({ page, request }) => {
  await addRun('saved-history', 'Completed', 1794225600);
  await publishLog('saved-history', 'Last timestamp', { timestamp: '2026-11-09T12:00:02Z', logLevel: 'ERROR' });
  await publishLog('saved-history', 'Repeated output');
  await publishLog('saved-history', 'Repeated output');
  const url = `${process.env.LOGCOLLECTOR_URL}/app/cicd/logs/saved-history`;
  await expect.poll(async () => (await (await request.get(url)).json()).data?.length).toBe(3);
  const read = async () => (await (await request.get(url)).json()).data;
  const records = await read();
  expect(records.map((record: { message: string }) => record.message).sort())
    .toEqual(['Last timestamp', 'Repeated output', 'Repeated output']);
  // Display the numeric severity actually persisted by the unchanged backend.
  expect(records.find((record: { message: string }) => record.message === 'Last timestamp').logLevel).toBe(0);
  expect(new Set(records.map((record: { id: string }) => record.id)).size).toBe(3);
  for (let index = 0; index < 3; index++) {
    expect((await read()).map((record: { id: string }) => record.id).sort())
      .toEqual(records.map((record: { id: string }) => record.id).sort());
  }
  // Known historical IDs make the equal-timestamp presentation order observable as text.
  await database.query(`INSERT INTO logcollector.log_entries (id, run_id, timestamp, log_level, message)
    VALUES ($1, $2, $3, $4, $5), ($6, $2, $3, $4, $7)`,
    ['z-tie', 'saved-history', '2026-11-09T12:00:01Z', 0, 'Tie second', 'a-tie', 'Tie first']);
  await page.goto('/runs/saved-history');
  const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
  await expect(logs.getByRole('listitem')).toHaveCount(5);
  await expect(logs.getByRole('listitem')).toHaveText([
    '2026-11-09T12:00:00Z INFO Repeated output',
    '2026-11-09T12:00:00Z INFO Repeated output',
    '2026-11-09T12:00:01Z INFO Tie first',
    '2026-11-09T12:00:01Z INFO Tie second',
    '2026-11-09T12:00:02Z INFO Last timestamp',
  ]);
  await logs.getByRole('button', { name: 'Refresh saved logs' }).click();
  await expect(logs.getByRole('button', { name: 'Refresh saved logs' })).toBeEnabled();
  await expect(logs.getByRole('listitem')).toHaveCount(5);
  await expect(logs.getByRole('status')).toHaveText('5 saved records.');
});

test('saved logs distinguish loading, active waiting, completed empty and unavailable snapshots', async ({ page, request }) => {
  await addRun('waiting-output', 'Running', 1794225600);
  await addRun('empty-history', 'Completed', 1794225500);
  let release!: () => void;
  const gate = new Promise<void>((resolve) => { release = resolve; });
  await page.route('**/api/logs/waiting-output', async (route) => {
    await gate;
    await route.continue();
  });
  const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
  try {
    await page.goto('/runs/waiting-output');
    await expect(logs).toContainText('Loading saved logs…');
    await expect(logs.getByRole('button', { name: 'Download logs', exact: true })).toBeDisabled();
    await expect(page.getByRole('region', { name: 'Run details', exact: true }).getByText('Running', { exact: true })).toBeVisible();
    release();
    await expect(logs).toContainText('Waiting for output');
  } finally {
    release();
    await page.unrouteAll({ behavior: 'wait' });
  }
  const empty = await request.get(`${process.env.LOGCOLLECTOR_URL}/app/cicd/logs/empty-history`);
  expect(await empty.json()).toMatchObject({ status: 200, error: '', data: [] });
  await publishLog('waiting-output', 'Active run already has output.');
  await waitForSavedRecords(request, 'waiting-output', 1);
  await logs.getByRole('button', { name: 'Refresh saved logs' }).click();
  await expect(logs).toContainText('Active run already has output.');
  await page.goto('/runs/empty-history');
  await expect(logs).toContainText('No saved logs for this run.');
  await expect(logs.getByRole('button', { name: 'Download logs', exact: true })).toBeEnabled();
  await logs.getByRole('searchbox', { name: 'Search logs' }).fill('no output');
  await expect(logs).not.toContainText('No saved log records match this search.');
  await expect(logs).not.toContainText('Waiting for output');
  service('stop', 'logcollector');
  try {
    await logs.getByRole('button', { name: 'Refresh saved logs' }).click();
    await expect(logs.getByRole('alert')).toContainText('stale');
    await expect(logs).not.toContainText('No saved logs for this run.');
    await page.reload();
    await expect(logs.getByRole('alert')).toContainText('log collector is unavailable');
    await expect(logs.getByRole('button', { name: 'Download logs', exact: true })).toBeDisabled();
    await expect(page.getByRole('region', { name: 'Run details', exact: true }).getByText('empty-history', { exact: true })).toBeVisible();
    service('start', 'logcollector');
    await waitForCollector(request);
    await logs.getByRole('button', { name: 'Refresh saved logs' }).click();
    await expect(logs).toContainText('No saved logs for this run.');
    await expect(logs.getByRole('alert')).toHaveCount(0);
  } finally {
    service('start', 'logcollector');
  }
});

test('a log-collector outage retains stale saved output beside usable run details and retry recovers', async ({ page, request }) => {
  await addRun('retained-output', 'Failed', 1794225600, 'cicd', { error: 'Compiler exited with code 2.' });
  await publishLog('retained-output', 'Previously saved compiler output.');
  await waitForSavedRecords(request, 'retained-output', 1);
  await page.clock.install();
  await page.goto('/runs/retained-output');
  const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
  await expect(logs).toContainText('Previously saved compiler output.');
  service('stop', 'logcollector');
  try {
    await logs.getByRole('button', { name: 'Refresh saved logs' }).click();
    await expect(logs.getByRole('alert')).toContainText('stale');
    await expect(logs).toContainText('Previously saved compiler output.');
    await expect(page.getByRole('heading', { name: 'Failure reason', exact: true })).toBeVisible();
    await expect(page.getByText('Compiler exited with code 2.', { exact: true })).toBeVisible();
    await page.clock.runFor(5000);
    await expect(page.getByRole('button', { name: 'Refresh runs' })).toBeEnabled();
    await expect(logs.getByRole('alert')).toContainText('stale');
    service('start', 'logcollector');
    await waitForCollector(request);
    await publishLog('retained-output', 'Collector recovered.', { timestamp: '2026-11-09T12:00:01Z' });
    await waitForSavedRecords(request, 'retained-output', 2);
    await logs.getByRole('button', { name: 'Refresh saved logs' }).click();
    await expect(logs.getByRole('listitem')).toHaveCount(2);
    await expect(logs).toContainText('Collector recovered.');
    await expect(logs.getByRole('alert')).toHaveCount(0);
    service('restart', 'ui');
    await expect.poll(async () => request.get('/').then((response) => response.status()).catch(() => 0)).toBe(200);
    await page.reload();
    await expect(logs.getByRole('listitem')).toHaveCount(2);
  } finally {
    service('start', 'logcollector');
  }
});

test('switching runs cancels ownership of a delayed saved-log response', async ({ page, request }) => {
  for (const id of ['log-run-a', 'log-run-b']) {
    await addRun(id, 'Completed', 1794225600);
    await publishLog(id, `Output belongs to ${id}.`);
    await waitForSavedRecords(request, id, 1);
  }
  const delayed = await delaySavedLogResponse(page, 'log-run-a');
  const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
  try {
    await page.goto('/');
    await page.getByRole('link', { name: 'Open run log-run-a', exact: true }).click();
    await delayed.received;
    await expect(logs).toContainText('Loading saved logs…');
    await page.goBack();
    await page.getByRole('link', { name: 'Open run log-run-b', exact: true }).click();
    await expect(logs).toContainText('Output belongs to log-run-b.');
    await expect(logs).not.toContainText('Output belongs to log-run-a.');
    delayed.release();
    await delayed.completed;
    await expect(logs.getByRole('listitem')).toHaveCount(1);
    await expect(logs).toContainText('Output belongs to log-run-b.');
    await page.goBack();
    await page.getByRole('link', { name: 'Open run log-run-a', exact: true }).click();
    await expect(logs).toContainText('Output belongs to log-run-a.');
    await expect(logs).not.toContainText('Output belongs to log-run-b.');
  } finally {
    delayed.release();
  }
});

test('saved logs preserve legacy records and safe Unicode content through encoded IDs and keyboard navigation', async ({ page, request }, testInfo) => {
  const id = 'build-δοκιμή #1?100%';
  const message = `Καλημέρα 👻 <img src=x onerror=alert(1)> <script>alert(1)</script>\n${'very-long-line-'.repeat(100)}\n${'compiler output\n'.repeat(80)}`;
  await addRun(id, 'Completed', 1794225600);
  // A historical row using exactly the pre-existing schema, without timestamp metadata.
  await database.query('INSERT INTO logcollector.log_entries (id, run_id, timestamp, log_level, message) VALUES ($1, $2, $3, $4, $5)',
    ['legacy-id', id, '', 0, 'Legacy saved output.']);
  await publishLog(id, message, { logLevel: 'WARN' });
  await publishLog(id, 'Debug output.', { timestamp: '2026-11-09T12:00:01Z', logLevel: 'DEBUG' });
  await waitForSavedRecords(request, id, 3);
  const mutations: string[] = [];
  const errors: string[] = [];
  page.on('request', (req) => { if (!['GET', 'HEAD'].includes(req.method())) mutations.push(req.method()); });
  page.on('pageerror', (error) => errors.push(error.message));
  await page.goto(`/runs/${encodeURIComponent(id)}`);
  const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
  await expect(logs.getByRole('listitem')).toHaveCount(3);
  await expect(logs.getByRole('listitem').first()).toHaveText('Timestamp unavailable INFO Legacy saved output.');
  await expect(logs).toContainText('Live connection open.');
  await expect(logs.getByText(message, { exact: true })).toBeVisible();
  await expect(logs.getByRole('listitem').nth(1)).toContainText('INFO');
  await expect(logs.getByRole('listitem').last()).toHaveText('2026-11-09T12:00:01Z INFO Debug output.');
  await expect(page.locator('img, #saved-logs script')).toHaveCount(0);
  for (const width of [1280, 390]) {
    await page.setViewportSize({ width, height: 900 });
    await logs.scrollIntoViewIfNeeded();
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
    await page.screenshot({ path: testInfo.outputPath(`saved-logs-${width}.png`), fullPage: true });
  }
  const output = logs.getByRole('list', { name: 'Saved log records' });
  await logs.getByRole('button', { name: 'Refresh saved logs' }).focus();
  await page.keyboard.press('Tab');
  await expect(logs.getByRole('searchbox', { name: 'Search logs' })).toBeFocused();
  await page.keyboard.press('Tab');
  await expect(logs.getByRole('button', { name: 'Download logs', exact: true })).toBeFocused();
  await page.keyboard.press('Tab');
  await expect(output).toBeFocused();
  await page.keyboard.press('End');
  await expect.poll(async () => output.evaluate((node) => node.scrollTop)).toBeGreaterThan(0);
  expect(mutations).toEqual([]);
  expect(errors).toEqual([]);
  expect((await request.post(`/api/logs/${encodeURIComponent(id)}`, { data: {} })).status()).toBe(403);
  expect((await request.get('/api/logs/ws/other-run')).status()).toBe(400);
  expect((await request.post('/api/logs/ws/other-run', { data: {} })).status()).toBe(403);
});
