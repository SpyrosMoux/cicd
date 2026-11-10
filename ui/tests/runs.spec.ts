import { test, expect } from '@playwright/test';
import { addRun, database, publishLog, resetRuns, service } from './fixtures';

test.beforeEach(resetRuns);
test.afterAll(() => database.end());

test('operator sees stored runs newest first with a deterministic tie-breaker', async ({ page }) => {
  await addRun('older', 'Completed', 1794222000, 'cicd');
  await addRun('b-newest', 'Running', 1794225600);
  await addRun('a-newest', 'Pending', 1794225600, 'notes');
  await page.goto('/');
  const rows = page.getByRole('row').filter({ has: page.getByRole('cell') });
  await expect(rows).toHaveCount(3);
  await expect(rows.nth(0)).toContainText('a-newest');
  await expect(rows.nth(1)).toContainText('b-newest');
  await expect(rows.nth(2)).toContainText('older');
  await expect(rows.nth(0)).toContainText('notes');
  await expect(rows.nth(1)).toContainText('recipe-book');
  await expect(rows.nth(2)).toContainText('cicd');
  await expect(rows.nth(0)).toContainText('main');
  await expect(rows.nth(0)).toContainText('Push');
  await expect(rows.nth(0)).toContainText('SpyrosMoux');
});

test('operator sees explicit statuses and honest lifecycle timing', async ({ page }) => {
  for (const [index, status] of ['Pending', 'Running', 'Completed', 'Failed', 'Canceled', 'Paused'].entries()) {
    await addRun(`run-${status}`, status, 1794225600 - index, 'cicd',
      status === 'Completed' ? { time_started: 1794225610, time_ended: 1794225675 } : {});
  }
  await addRun('missing-timing', 'Pending', 0);
  await page.goto('/');
  const row = (id: string) => page.getByRole('row').filter({ hasText: id });
  for (const status of ['Pending', 'Running', 'Completed', 'Failed', 'Canceled']) {
    await expect(row(`run-${status}`).getByText(status, { exact: true })).toBeVisible();
  }
  await expect(row('run-Paused')).toContainText('Unknown (Paused)');
  await expect(row('run-Completed')).toContainText('Started: 9 Nov 2026, 12:00:10 UTC');
  await expect(row('run-Completed')).toContainText('Ended: 9 Nov 2026, 12:01:15 UTC');
  await expect(row('run-Completed')).toContainText('Duration: 1m 5s');
  await expect(row('run-Pending')).toContainText('Started: Pending');
  await expect(row('run-Running')).toContainText('Duration: Unavailable');
  await expect(row('missing-timing')).toContainText('Triggered: Unavailable');
  await expect(page.getByText(/1970/)).toHaveCount(0);
});

test('operator can distinguish initial loading from an empty installation', async ({ page }) => {
  let release!: () => void;
  const gate = new Promise<void>((resolve) => { release = resolve; });
  // Delay, then forward the real request. No production response is mocked.
  await page.route('**/api/runs', async (route) => { await gate; await route.continue(); });
  try {
    await page.goto('/');
    await expect(page.getByRole('status')).toHaveText('Loading runs…');
    await expect(page.getByRole('heading', { name: 'No runs yet' })).toHaveCount(0);
  } finally {
    release();
  }
  await expect(page.getByRole('heading', { name: 'No runs yet' })).toBeVisible();
  await expect(page.getByText('0 runs', { exact: true })).toBeVisible();
});

test('operator can recover from an API outage without mistaking it for empty history', async ({ page, request }) => {
  await addRun('retained-run', 'Completed', 1794225600);
  service('stop', 'api');
  try {
    await page.goto('/');
    await expect(page.getByRole('heading', { name: 'Unable to load runs' })).toBeVisible();
    await expect(page.getByRole('heading', { name: 'No runs yet' })).toHaveCount(0);
    service('start', 'api');
    await expect.poll(async () => (await request.get('/api/runs')).status()).toBe(200);
    await page.getByRole('button', { name: 'Try again', exact: true }).click();
    await expect(page.getByRole('row').filter({ hasText: 'retained-run' })).toBeVisible();
    service('stop', 'api');
    await page.getByRole('button', { name: 'Refresh runs' }).click();
    await expect(page.getByRole('alert')).toContainText('stale');
    await expect(page.getByRole('row').filter({ hasText: 'retained-run' })).toBeVisible();
  } finally {
    service('start', 'api');
  }
});

test('operator retains access to backend history after a UI restart', async ({ page, request }) => {
  await addRun('persisted-run', 'Completed', 1794225600, 'cicd');
  await page.goto('/');
  await expect(page.getByRole('row').filter({ hasText: 'persisted-run' })).toBeVisible();
  service('restart', 'ui');
  await expect.poll(async () => {
    try { return (await request.get('/api/runs')).status(); } catch { return 0; }
  }).toBe(200);
  await page.reload();
  await expect(page.getByRole('row').filter({ hasText: 'persisted-run' })).toBeVisible();
});

test('operator can read safe text on desktop and narrow screens using the keyboard', async ({ page }, testInfo) => {
  const repository = 'recipes Συνταγές <img src=x onerror=alert(1)>';
  const branch = `feature/${'long-branch-'.repeat(40)}`;
  await addRun('safe-run', 'Completed', 1794225600, repository, { branch });
  const mutations: string[] = [];
  const errors: string[] = [];
  page.on('request', (request) => {
    if (!['GET', 'HEAD'].includes(request.method())) mutations.push(request.method());
  });
  page.on('pageerror', (error) => errors.push(error.message));
  await page.goto('/');
  await expect(page.getByRole('cell').getByText(repository, { exact: true })).toBeVisible();
  await expect(page.getByText(branch, { exact: true })).toBeVisible();
  await expect(page.locator('img')).toHaveCount(0);
  await page.keyboard.press('Tab');
  await expect(page.getByRole('link', { name: 'FlowForge home' })).toBeFocused();
  await page.keyboard.press('Tab');
  await expect(page.getByRole('button', { name: 'Refresh runs' })).toBeFocused();
  await Promise.all([
    page.waitForResponse((response) => response.url().endsWith('/api/runs')),
    page.keyboard.press('Enter'),
  ]);
  await expect(page.getByRole('button', { name: 'Refresh runs' })).toBeEnabled();
  await page.getByRole('button', { name: 'Refresh runs' }).focus();
  await page.keyboard.press('Tab');
  await expect(page.getByRole('combobox', { name: 'Repository', exact: true })).toBeFocused();
  await page.keyboard.press('r');
  await expect(page.getByRole('combobox', { name: 'Repository', exact: true })).toHaveValue(repository);
  await page.keyboard.press('Tab');
  await expect(page.getByRole('combobox', { name: 'Status', exact: true })).toBeFocused();
  await page.keyboard.press('p');
  await expect(page.getByRole('heading', { name: 'No matching runs' })).toBeVisible();
  await page.keyboard.press('Tab');
  await expect(page.getByRole('button', { name: 'Clear filters' })).toBeFocused();
  await page.keyboard.press('Enter');
  await expect(page.getByRole('cell').getByText(repository, { exact: true })).toBeVisible();
  for (const width of [1280, 390]) {
    await page.setViewportSize({ width, height: 900 });
    await expect(page.getByRole('cell').getByText(repository, { exact: true })).toBeVisible();
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
    await page.screenshot({ path: testInfo.outputPath(`runs-${width}.png`), fullPage: true });
  }
  expect(mutations).toEqual([]);
  expect(errors).toEqual([]);
});

test('UI routing exposes run reads without forwarding mutation operations', async ({ request }) => {
  expect((await request.get('/api/runs')).status()).toBe(200);
  expect((await request.post('/api/runs', { data: {} })).status()).toBe(403);
  expect((await request.put('/api/runs/example', { data: { status: 'Failed' } })).status()).toBe(404);
  expect((await request.post('/app/cicd/api/gh/webhook', { data: {} })).status()).toBe(405);
});

test('shared log fixtures pass through RabbitMQ and the real log collector', async ({ request }) => {
  await publishLog('fixture-run', 'Fixture output: καλημέρα <build>');
  await expect.poll(async () => {
    const response = await request.get(`${process.env.LOGCOLLECTOR_URL}/app/cicd/logs/fixture-run`);
    return ((await response.json()).data ?? []).map((entry: { message: string }) => entry.message);
  }).toEqual(['Fixture output: καλημέρα <build>']);
});

test('operator combines repository and status filters, distinguishes no matches, and clears them', async ({ page }) => {
  await addRun('cicd-completed', 'Completed', 1794225603, 'cicd');
  await addRun('cicd-running', 'Running', 1794225602, 'cicd');
  await addRun('notes-failed', 'Failed', 1794225601, 'notes');
  await page.goto('/');
  const repository = page.getByRole('combobox', { name: 'Repository', exact: true });
  const status = page.getByRole('combobox', { name: 'Status', exact: true });
  const rows = page.getByRole('row').filter({ has: page.getByRole('cell') });
  await expect(rows).toHaveCount(3);
  await expect(status.locator('option')).toHaveText([
    'All statuses', 'Pending', 'Running', 'Completed', 'Failed', 'Canceled',
  ]);
  await repository.selectOption('cicd');
  await expect(rows).toHaveCount(2);
  await status.selectOption('Running');
  await expect(rows).toHaveCount(1);
  await expect(rows).toContainText('cicd-running');
  await repository.selectOption('');
  await status.selectOption('Failed');
  await expect(rows).toHaveCount(1);
  await expect(rows).toContainText('notes-failed');
  await repository.selectOption('cicd');
  await expect(page.getByRole('heading', { name: 'No matching runs' })).toBeVisible();
  await expect(page.getByRole('heading', { name: 'No runs yet' })).toHaveCount(0);
  await expect(page.getByText('0 of 3 runs', { exact: true })).toBeVisible();
  const clear = page.getByRole('button', { name: 'Clear filters' });
  await clear.focus();
  await page.keyboard.press('Enter');
  await expect(repository).toHaveValue('');
  await expect(status).toHaveValue('');
  await expect(rows).toHaveCount(3);
  await expect(clear).toBeDisabled();
});

test('operator automatically sees new runs and status changes without losing active filters', async ({ page, request }) => {
  await addRun('active-run', 'Running', 1794225600, 'cicd');
  await page.clock.install();
  await page.goto('/');
  const repository = page.getByRole('combobox', { name: 'Repository', exact: true });
  const status = page.getByRole('combobox', { name: 'Status', exact: true });
  const rows = page.getByRole('row').filter({ has: page.getByRole('cell') });
  await expect(rows).toHaveCount(1);
  await repository.selectOption('cicd');
  await status.selectOption('Running');
  await status.focus();
  await addRun('b-arrival', 'Running', 1794225700, 'cicd');
  await addRun('a-arrival', 'Running', 1794225700, 'cicd');
  await addRun('other-project', 'Failed', 1794225701, 'new-project');
  const response = await request.put(
    `http://127.0.0.1:${process.env.API_TEST_PORT}/app/cicd/api/runs/active-run`,
    { data: { status: 'Completed', time_ended: 1794225700 } },
  );
  expect(response.ok()).toBe(true);
  await page.clock.runFor(5000);
  await expect(rows).toHaveCount(2);
  await expect(rows.nth(0)).toContainText('a-arrival');
  await expect(rows.nth(1)).toContainText('b-arrival');
  await expect(repository).toHaveValue('cicd');
  await expect(status).toHaveValue('Running');
  await expect(status).toBeFocused();
  await expect(repository.locator('option')).toHaveText(['All repositories', 'cicd', 'new-project']);
  await status.selectOption('Completed');
  await expect(rows).toHaveCount(1);
  await expect(rows).toContainText('active-run');
  await expect(rows).toContainText('Completed');
  await repository.selectOption('new-project');
  await status.selectOption('Failed');
  await expect(rows).toContainText('other-project');
});

test('operator keeps their reading position and selected text when new runs arrive', async ({ page }) => {
  for (let index = 0; index < 30; index++) {
    await addRun(`history-${index}`, 'Completed', 1794225600 - index, 'cicd');
  }
  await page.clock.install();
  await page.goto('/');
  const anchor = page.getByRole('row').filter({ hasText: 'history-15' });
  await anchor.scrollIntoViewIfNeeded();
  const top = (await anchor.boundingBox())!.y;
  await anchor.getByText('history-15', { exact: true }).evaluate((node) => {
    const selection = window.getSelection()!;
    const range = document.createRange();
    range.selectNodeContents(node);
    selection.removeAllRanges();
    selection.addRange(range);
  });
  await addRun('new-arrival', 'Pending', 1794225700, 'cicd');
  await page.clock.runFor(5000);
  await expect(page.getByRole('row').filter({ hasText: 'new-arrival' })).toHaveCount(1);
  expect(Math.abs((await anchor.boundingBox())!.y - top)).toBeLessThan(2);
  expect(await page.evaluate(() => window.getSelection()?.toString())).toBe('history-15');
});

test('automatic refresh is bounded and never overlaps a slow request', async ({ page }) => {
  await addRun('retained-run', 'Running', 1794225600, 'cicd');
  await page.clock.install();
  let requests = 0;
  let release!: () => void;
  const gate = new Promise<void>((resolve) => { release = resolve; });
  await page.route('**/api/runs', async (route) => {
    requests++;
    if (requests === 2) await gate;
    await route.continue();
  });
  try {
    await page.goto('/');
    const refresh = page.getByRole('button', { name: 'Refresh runs' });
    await expect(refresh).toBeEnabled();
    expect(requests).toBe(1);
    await page.clock.runFor(4999);
    expect(requests).toBe(1);
    await page.clock.runFor(1);
    await expect.poll(() => requests).toBe(2);
    await page.clock.runFor(6000);
    expect(requests).toBe(2);
    await expect(refresh).toBeDisabled();
    await expect(page.getByRole('row').filter({ hasText: 'retained-run' })).toBeVisible();
    release();
    await expect(refresh).toBeEnabled();
    await page.clock.runFor(4999);
    expect(requests).toBe(2);
    await page.clock.runFor(1);
    await expect.poll(() => requests).toBe(3);
    await expect(refresh).toBeEnabled();
    await refresh.click();
    await expect.poll(() => requests).toBe(4);
    await expect(refresh).toBeEnabled();
    await page.clock.runFor(4999);
    expect(requests).toBe(4);
  } finally {
    release();
  }
});

test('operator retains filtered history through an outage and recovers automatically', async ({ page, request }) => {
  await addRun('retained-run', 'Running', 1794225600, 'cicd');
  await page.clock.install();
  await page.goto('/');
  const repository = page.getByRole('combobox', { name: 'Repository', exact: true });
  const status = page.getByRole('combobox', { name: 'Status', exact: true });
  await expect(page.getByRole('row').filter({ hasText: 'retained-run' })).toBeVisible();
  await repository.selectOption('cicd');
  await status.selectOption('Running');
  service('stop', 'api');
  try {
    await page.clock.runFor(5000);
    await expect(page.getByRole('alert')).toContainText('stale');
    await expect(page.getByRole('row').filter({ hasText: 'retained-run' })).toBeVisible();
    await expect(repository).toHaveValue('cicd');
    await expect(status).toHaveValue('Running');
    await expect(page.getByRole('heading', { name: 'No runs yet' })).toHaveCount(0);
    service('start', 'api');
    await expect.poll(async () => (await request.get('/api/runs')).status()).toBe(200);
    await addRun('recovered-run', 'Running', 1794225700, 'cicd');
    await page.clock.runFor(5000);
    await expect(page.getByRole('row').filter({ hasText: 'recovered-run' })).toBeVisible();
    await expect(page.getByRole('alert')).toHaveCount(0);
    await expect(repository).toHaveValue('cicd');
    await expect(status).toHaveValue('Running');
    // Even if the selected repository leaves the latest snapshot, it stays selected.
    await resetRuns();
    await addRun('different-project', 'Running', 1794225800, 'notes');
    await page.clock.runFor(5000);
    await expect(page.getByRole('heading', { name: 'No matching runs' })).toBeVisible();
    await expect(repository).toHaveValue('cicd');
    await expect(repository.locator('option')).toHaveText(['All repositories', 'notes', 'cicd']);
    await page.getByRole('button', { name: 'Clear filters' }).click();
    await expect(page.getByRole('row').filter({ hasText: 'different-project' })).toBeVisible();
  } finally {
    service('start', 'api');
  }
});
