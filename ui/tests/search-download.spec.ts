import { readFile } from 'node:fs/promises';
import { expect } from '@playwright/test';
import { addRun, database, publishLog, resetRuns, service, test, waitForCollector, waitForSavedRecords } from './fixtures';

test.beforeEach(resetRuns);

test('log search matches literal text and Unicode, preserves repeated records and clears no matches', async ({ page, request }) => {
  await page.routeWebSocket('**/api/logs/ws/**', (socket) => socket.close());
  await addRun('search-history', 'Completed', 1794225600);
  for (const message of ['Repeated [build.*] output', 'Repeated [build.*] output', 'buildXYZ output', 'Καλημέρα 👻 <img src=x onerror=alert(1)>']) {
    await publishLog('search-history', message);
  }
  await waitForSavedRecords(request, 'search-history', 4);
  await page.goto('/runs/search-history');
  const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
  const search = logs.getByRole('searchbox', { name: 'Search logs', exact: true });
  const clear = logs.getByRole('button', { name: 'Clear log search', exact: true });
  const lines = logs.getByRole('listitem');
  await expect(lines).toHaveCount(4);
  await expect(clear).toBeDisabled();
  await search.fill('2026-11-09T12:00:00Z INFO');
  await expect(lines).toHaveCount(4);
  await search.fill('repeated');
  await expect(lines).toHaveCount(0);
  await search.fill('[build.*]');
  await expect(lines).toHaveCount(2);
  await expect(lines).toHaveText([
    '2026-11-09T12:00:00Z INFO Repeated [build.*] output',
    '2026-11-09T12:00:00Z INFO Repeated [build.*] output',
  ]);
  await expect(logs).toContainText('2 of 4 saved records match this search.');
  await search.fill('Καλημέρα 👻');
  await expect(lines).toHaveCount(1);
  await expect(lines).toContainText('Καλημέρα 👻 <img src=x onerror=alert(1)>');
  await search.fill('<img src=x onerror=alert(1)>');
  await expect(lines).toHaveCount(1);
  await expect(logs.locator('img')).toHaveCount(0);
  await search.fill('no such output');
  await expect(lines).toHaveCount(0);
  await expect(logs).toContainText('No saved log records match this search.');
  await expect(logs).not.toContainText('No saved logs for this run.');
  await clear.click();
  await expect(search).toHaveValue('');
  await expect(search).toBeFocused();
  await expect(lines).toHaveCount(4);
  await expect(logs).not.toContainText('No saved log records match this search.');
  await expect(clear).toBeDisabled();
  await expect.poll(async () => (await (await request.get(`${process.env.LOGCOLLECTOR_URL}/app/cicd/logs/search-history`)).json()).data?.length).toBe(4);
});

test('log search and export use updated saved snapshots without a live connection', async ({ page, request }) => {
  await page.routeWebSocket('**/api/logs/ws/**', (socket) => socket.close());
  await addRun('updated-search', 'Completed', 1794225600);
  await publishLog('updated-search', 'needle first');
  await publishLog('updated-search', 'hidden first');
  await waitForSavedRecords(request, 'updated-search', 2);
  await page.clock.install();
  await page.goto('/runs/updated-search');
  const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
  const search = logs.getByRole('searchbox', { name: 'Search logs' });
  await expect(logs.getByRole('listitem')).toHaveCount(2);
  await search.fill('needle');
  await expect(logs.getByRole('listitem')).toHaveCount(1);
  await publishLog('updated-search', 'needle later', { timestamp: '2026-11-09T12:00:01Z' });
  await publishLog('updated-search', 'hidden later', { timestamp: '2026-11-09T12:00:02Z' });
  await waitForSavedRecords(request, 'updated-search', 4);
  await page.clock.runFor(5000);
  await expect(logs.getByRole('listitem')).toHaveCount(2);
  await expect(search).toHaveValue('needle');
  await expect(logs).toContainText('2 of 4 saved records match this search.');
  const downloading = page.waitForEvent('download');
  await logs.getByRole('button', { name: 'Download logs', exact: true }).click();
  const text = await readFile((await (await downloading).path())!, 'utf8');
  for (const message of ['needle first', 'hidden first', 'needle later', 'hidden later']) {
    expect(text).toContain(message);
  }
  await logs.getByRole('button', { name: 'Clear log search' }).click();
  await expect(logs.getByRole('listitem')).toHaveCount(4);
});

test('log downloads stay with the selected run when an earlier snapshot arrives late', async ({ page, request }) => {
  for (const id of ['export-a', 'export-b']) {
    await addRun(id, 'Completed', 1794225600);
    await publishLog(id, `Output only for ${id}.`);
    await waitForSavedRecords(request, id, 1);
  }
  let release!: () => void;
  let received!: () => void;
  let finished!: () => void;
  const gate = new Promise<void>((resolve) => { release = resolve; });
  const waiting = new Promise<void>((resolve) => { received = resolve; });
  const completed = new Promise<void>((resolve) => { finished = resolve; });
  let delayed = false;
  await page.route('**/api/logs/export-a', async (route) => {
    if (delayed) { await route.continue(); return; }
    delayed = true;
    const response = await route.fetch();
    received();
    await gate;
    try { await route.fulfill({ response }); } finally { finished(); }
  });
  const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
  const downloadButton = logs.getByRole('button', { name: 'Download logs', exact: true });
  const search = logs.getByRole('searchbox', { name: 'Search logs' });
  try {
    await page.goto('/');
    await page.getByRole('link', { name: 'Open run export-a', exact: true }).click();
    await waiting;
    await expect(downloadButton).toBeDisabled();
    await search.fill('only for export-a');
    await page.goBack();
    await page.getByRole('link', { name: 'Open run export-b', exact: true }).click();
    await expect(search).toHaveValue('');
    await expect(logs.getByRole('listitem')).toHaveText(['2026-11-09T12:00:00Z INFO Output only for export-b.']);
    await search.fill('no matching lines');
    release();
    await completed;
    const downloading = page.waitForEvent('download');
    await downloadButton.click();
    const download = await downloading;
    expect(download.suggestedFilename()).toBe('flowforge-export-b-logs.txt');
    const text = await readFile((await download.path())!, 'utf8');
    expect(text).toContain('Run: export-b\n');
    expect(text).toContain('Output only for export-b.');
    expect(text).not.toContain('export-a');
    await expect(search).toHaveValue('no matching lines');
    await expect(logs.getByRole('listitem')).toHaveCount(0);
  } finally {
    release();
  }
});

test('log searches and downloads retain honest stale output during an outage and recover', async ({ page, request }) => {
  await addRun('stale-export', 'Failed', 1794225600);
  await publishLog('stale-export', 'Retained compiler output.');
  await waitForSavedRecords(request, 'stale-export', 1);
  await page.goto('/runs/stale-export');
  const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
  await expect(logs.getByRole('listitem')).toHaveCount(1);
  const search = logs.getByRole('searchbox', { name: 'Search logs' });
  await search.fill('absent');
  service('stop', 'logcollector');
  try {
    await logs.getByRole('button', { name: 'Refresh saved logs' }).click();
    await expect(logs.getByRole('alert')).toContainText('stale');
    await logs.getByRole('button', { name: 'Clear log search' }).click();
    await expect(logs.getByRole('listitem')).toHaveCount(1);
    await expect(logs.getByRole('alert')).toContainText('stale');
    await search.fill('absent');
    await expect(logs).toContainText('No saved log records match this search.');
    const downloading = page.waitForEvent('download');
    await logs.getByRole('button', { name: 'Download logs', exact: true }).click();
    const text = await readFile((await (await downloading).path())!, 'utf8');
    expect(text).toContain('Status: Retained output; last saved-log refresh failed.');
    expect(text).toContain('Run: stale-export\n');
    expect(text).toContain('2026-11-09T12:00:00Z INFO Retained compiler output.');
    service('start', 'logcollector');
    await waitForCollector(request);
    await logs.getByRole('button', { name: 'Refresh saved logs' }).click();
    await expect(logs.getByRole('alert')).toHaveCount(0);
    await expect(search).toHaveValue('absent');
    await expect(logs).toContainText('No saved log records match this search.');
  } finally {
    service('start', 'logcollector');
  }
});

test('log search, clearing and unfiltered download work with the keyboard at narrow widths', async ({ page, request }, testInfo) => {
  await addRun('keyboard-export', 'Completed', 1794225600);
  await publishLog('keyboard-export', 'Keyboard-accessible output 👻');
  await waitForSavedRecords(request, 'keyboard-export', 1);
  const mutations: string[] = [];
  page.on('request', (req) => { if (!['GET', 'HEAD'].includes(req.method())) mutations.push(req.method()); });
  await page.goto('/runs/keyboard-export');
  const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
  const search = logs.getByRole('searchbox', { name: 'Search logs' });
  await expect(logs.getByRole('listitem')).toHaveCount(1);
  await logs.getByRole('button', { name: 'Refresh saved logs' }).focus();
  await page.keyboard.press('Tab');
  await expect(search).toBeFocused();
  await page.keyboard.type('absent');
  await expect(logs.getByRole('listitem')).toHaveCount(0);
  await page.keyboard.press('Tab');
  await expect(logs.getByRole('button', { name: 'Clear log search' })).toBeFocused();
  await page.keyboard.press('Enter');
  await expect(search).toBeFocused();
  await expect(search).toHaveValue('');
  await expect(logs.getByRole('listitem')).toHaveCount(1);
  await page.keyboard.type('absent');
  await page.keyboard.press('Tab');
  await page.keyboard.press('Tab');
  await expect(logs.getByRole('button', { name: 'Download logs', exact: true })).toBeFocused();
  const downloading = page.waitForEvent('download');
  await page.keyboard.press('Enter');
  const text = await readFile((await (await downloading).path())!, 'utf8');
  expect(text).toContain('Keyboard-accessible output 👻');
  const noMatches = logs.getByText('No saved log records match this search.', { exact: false });
  await expect(noMatches).toBeVisible();
  await expect(noMatches).toHaveAttribute('aria-live', 'polite');
  for (const width of [1280, 390]) {
    await page.setViewportSize({ width, height: 900 });
    await logs.scrollIntoViewIfNeeded();
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
    await page.screenshot({ path: testInfo.outputPath(`search-download-${width}.png`), fullPage: true });
  }
  expect(mutations).toEqual([]);
});

test('log downloads export the unfiltered saved snapshot with run identity and readable metadata', async ({ page, request }) => {
  const id = 'build-δοκιμή #1?100%';
  await addRun(id, 'Completed', 1794225600);
  await addRun('other-export', 'Completed', 1794225500);
  await publishLog('other-export', 'Other run secret output.');
  await publishLog(id, 'Repeated output');
  await publishLog(id, 'Repeated output');
  await database.query(`INSERT INTO logcollector.log_entries (id, run_id, timestamp, log_level, message)
    VALUES ($1, $2, $3, $4, $5), ($6, $2, $7, $8, $9)`,
    ['legacy', id, '', 1, 'Καλημέρα 👻 <script>literal</script>', 'error', '2026-11-09T12:00:01Z', 2, 'Compiler failed.\nDetails on another line.']);
  await waitForSavedRecords(request, id, 4);
  await page.goto(`/runs/${encodeURIComponent(id)}`);
  const logs = page.getByRole('region', { name: 'Saved logs', exact: true });
  await expect(logs.getByRole('listitem')).toHaveCount(4);
  await logs.getByRole('searchbox', { name: 'Search logs' }).fill('WARN');
  await expect(logs.getByRole('listitem')).toHaveCount(1);
  const downloading = page.waitForEvent('download');
  await logs.getByRole('button', { name: 'Download logs', exact: true }).click();
  const download = await downloading;
  expect(download.suggestedFilename()).toBe('flowforge-build-δοκιμή__1_100_-logs.txt');
  const text = await readFile((await download.path())!, 'utf8');
  expect(text).toContain(`Run: ${id}\n`);
  expect(text).toMatch(/Snapshot loaded: \d{4}-\d{2}-\d{2}T/);
  expect(text).toContain('may be incomplete or stale');
  expect(text).toContain([
    'Timestamp unavailable WARN Καλημέρα 👻 <script>literal</script>',
    '2026-11-09T12:00:00Z INFO Repeated output',
    '2026-11-09T12:00:00Z INFO Repeated output',
    '2026-11-09T12:00:01Z ERROR Compiler failed.',
    'Details on another line.',
    '',
  ].join('\n'));
  expect(text).not.toContain('Other run secret output.');
  await expect(logs.getByRole('searchbox', { name: 'Search logs' })).toHaveValue('WARN');
  await expect(logs.getByRole('listitem')).toHaveCount(1);
});
