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
  const repository = 'Συνταγές <img src=x onerror=alert(1)>';
  const branch = `feature/${'long-branch-'.repeat(40)}`;
  await addRun('safe-run', 'Completed', 1794225600, repository, { branch });
  const mutations: string[] = [];
  const errors: string[] = [];
  page.on('request', (request) => {
    if (!['GET', 'HEAD'].includes(request.method())) mutations.push(request.method());
  });
  page.on('pageerror', (error) => errors.push(error.message));
  await page.goto('/');
  await expect(page.getByText(repository, { exact: true })).toBeVisible();
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
  for (const width of [1280, 390]) {
    await page.setViewportSize({ width, height: 900 });
    await expect(page.getByText(repository, { exact: true })).toBeVisible();
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
