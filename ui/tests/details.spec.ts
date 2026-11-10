import { expect } from '@playwright/test';
import { addRun, resetRuns, service, test } from './fixtures';

test.beforeEach(resetRuns);

test('operator opens a stable run detail link and keeps list context through browser navigation', async ({ page }) => {
  await addRun('selected-run', 'Completed', 1794225600, 'cicd', {
    branch: 'feature/details', trigger_type: 'Manual', triggered_by: 'operator',
    time_started: 1794225610, time_ended: 1794225675,
  });
  await addRun('other-run', 'Pending', 1794225500, 'notes');
  await page.goto('/');
  await page.getByRole('combobox', { name: 'Repository', exact: true }).selectOption('cicd');
  await page.getByRole('combobox', { name: 'Status', exact: true }).selectOption('Completed');
  const link = page.getByRole('link', { name: 'Open run selected-run', exact: true });
  await expect(link).toHaveAttribute('href', '/runs/selected-run');
  await link.focus();
  await page.keyboard.press('Enter');
  await expect(page).toHaveURL(/\/runs\/selected-run$/);
  await expect(page.getByRole('heading', { name: 'Run details', exact: true })).toBeVisible();
  const details = page.getByRole('region', { name: 'Run details', exact: true });
  for (const text of ['selected-run', 'cicd', 'feature/details', 'Completed', 'Manual', 'operator',
    '9 Nov 2026, 12:00:00 UTC', '9 Nov 2026, 12:00:10 UTC', '9 Nov 2026, 12:01:15 UTC', '1m 5s']) {
    await expect(details.getByText(text, { exact: true })).toBeVisible();
  }
  await page.goBack();
  await expect(page.getByRole('combobox', { name: 'Repository', exact: true })).toHaveValue('cicd');
  await expect(page.getByRole('combobox', { name: 'Status', exact: true })).toHaveValue('Completed');
  await expect(link).toBeFocused();
  await page.goForward();
  await expect(details).toContainText('selected-run');
  await page.reload();
  await expect(details).toContainText('selected-run');
  await page.getByRole('link', { name: 'Back to runs', exact: true }).click();
  await expect(page).toHaveURL(/\/$/);
  await expect(page.getByRole('combobox', { name: 'Repository', exact: true })).toHaveValue('cicd');
  await expect(page.getByRole('row').filter({ has: page.getByRole('cell') })).toHaveCount(1);
});

test('selected run remains readable as stale through an outage and recovers automatically without losing focus', async ({ page, request }) => {
  await addRun('retained-detail', 'Failed', 1794225600, 'cicd', { error: 'Compiler exited with code 2.' });
  await page.clock.install();
  await page.goto('/runs/retained-detail');
  const details = page.getByRole('region', { name: 'Run details', exact: true });
  await expect(details).toContainText('Compiler exited with code 2.');
  const back = page.getByRole('link', { name: 'Back to runs', exact: true });
  await back.focus();
  service('stop', 'api');
  try {
    await page.clock.runFor(5000);
    await expect(page.getByRole('alert')).toContainText('stale');
    await expect(details).toContainText('Compiler exited with code 2.');
    await expect(page.getByRole('heading', { name: 'Run not found', exact: true })).toHaveCount(0);
    await expect(back).toBeFocused();
    service('start', 'api');
    await expect.poll(async () => (await request.get('/api/runs')).status()).toBe(200);
    const response = await request.put(
      `http://127.0.0.1:${process.env.API_TEST_PORT}/app/cicd/api/runs/retained-detail`,
      { data: { status: 'Completed', time_started: 1794225610, time_ended: 1794225675 } },
    );
    expect(response.ok()).toBe(true);
    await page.clock.runFor(5000);
    await expect(details.getByText('Completed', { exact: true })).toBeVisible();
    await expect(details.getByText('1m 5s', { exact: true })).toBeVisible();
    await expect(page.getByRole('heading', { name: 'Failure reason', exact: true })).toHaveCount(0);
    await expect(page.getByRole('alert')).toHaveCount(0);
    await expect(back).toBeFocused();
    await expect(page).toHaveURL(/\/runs\/retained-detail$/);
  } finally {
    service('start', 'api');
  }
});

test('a delayed real response cannot overwrite a newly selected run', async ({ page }) => {
  await addRun('run-a', 'Completed', 1794225600, 'cicd');
  await addRun('run-b', 'Failed', 1794225500, 'notes', { error: 'Run B failed.' });
  await page.clock.install();
  let release!: () => void;
  let received!: () => void;
  const gate = new Promise<void>((resolve) => { release = resolve; });
  const waiting = new Promise<void>((resolve) => { received = resolve; });
  let requests = 0;
  await page.route('**/api/runs', async (route) => {
    if (++requests !== 2) { await route.continue(); return; }
    // Hold an actual API snapshot; its response body is not substituted.
    const response = await route.fetch();
    received();
    await gate;
    await route.fulfill({ response });
  });
  try {
    await page.goto('/');
    await page.getByRole('link', { name: 'Open run run-a', exact: true }).click();
    await page.clock.runFor(5000);
    await waiting;
    await page.goBack();
    await page.getByRole('link', { name: 'Open run run-b', exact: true }).click();
    const details = page.getByRole('region', { name: 'Run details', exact: true });
    await expect(details.getByText('run-b', { exact: true })).toBeVisible();
    release();
    await expect(page.getByRole('button', { name: 'Refresh runs' })).toBeEnabled();
    await expect(details.getByText('run-b', { exact: true })).toBeVisible();
    await expect(details).toContainText('Run B failed.');
    await expect(details.getByText('run-a', { exact: true })).toHaveCount(0);
    await expect(page).toHaveURL(/\/runs\/run-b$/);
    await page.goBack();
    await expect(page.getByRole('link', { name: 'Open run run-b', exact: true })).toBeFocused();
  } finally {
    release();
  }
});

test('returning from a refreshed detail link retains the list reading position', async ({ page }) => {
  for (let index = 0; index < 30; index++) {
    await addRun(`history-${index}`, 'Completed', 1794225600 - index, 'cicd');
  }
  await page.goto('/');
  const row = page.getByRole('row').filter({ hasText: 'history-15' });
  await row.scrollIntoViewIfNeeded();
  const top = (await row.boundingBox())!.y;
  await row.getByRole('link', { name: 'Open run history-15', exact: true }).click();
  let release!: () => void;
  const gate = new Promise<void>((resolve) => { release = resolve; });
  await page.route('**/api/runs', async (route) => { await gate; await route.continue(); });
  try {
    await page.reload();
    await expect(page.getByText('Loading run…', { exact: true })).toBeVisible();
    await page.getByRole('link', { name: 'Back to runs', exact: true }).click();
    release();
    await expect(row).toBeVisible();
    expect(Math.abs((await row.boundingBox())!.y - top)).toBeLessThan(2);
    await expect(row.getByRole('link', { name: 'Open run history-15', exact: true })).toBeFocused();
  } finally {
    release();
  }
});

test('run detail links and safe text remain accessible and read-only on desktop and narrow screens', async ({ page }, testInfo) => {
  const id = 'build/δοκιμή #1';
  const repository = 'cicd <img src=x onerror=alert(1)>';
  const branch = `feature/${'long-branch-'.repeat(40)}`;
  const error = `Compiler output: καλημέρα <script>alert(1)</script>\n${'long-output-'.repeat(60)}`;
  await addRun(id, 'Failed', 1794225600, repository, { branch, error });
  const mutations: string[] = [];
  const errors: string[] = [];
  page.on('request', (request) => {
    if (!['GET', 'HEAD'].includes(request.method())) mutations.push(request.method());
  });
  page.on('pageerror', (error) => errors.push(error.message));
  await page.goto(`/runs/${encodeURIComponent(id)}`);
  const details = page.getByRole('region', { name: 'Run details', exact: true });
  for (const text of [id, repository, branch, error]) {
    await expect(details.getByText(text, { exact: true })).toBeVisible();
  }
  await expect(page.locator('img')).toHaveCount(0);
  await expect(page.getByRole('button', { name: /trigger|retry|cancel|edit|delete/i })).toHaveCount(0);
  for (const width of [1280, 390]) {
    await page.setViewportSize({ width, height: 900 });
    await expect(details).toBeVisible();
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
    await page.screenshot({ path: testInfo.outputPath(`details-${width}.png`), fullPage: true });
  }
  await page.keyboard.press('Tab');
  await expect(page.getByRole('link', { name: 'FlowForge home' })).toBeFocused();
  await page.keyboard.press('Tab');
  await expect(page.getByRole('button', { name: 'Refresh runs' })).toBeFocused();
  await page.keyboard.press('Tab');
  await expect(page.getByRole('link', { name: 'Back to runs', exact: true })).toBeFocused();
  await page.keyboard.press('Enter');
  await expect(page.getByRole('heading', { name: 'Recent runs', exact: true })).toBeVisible();
  const link = page.getByRole('link', { name: `Open run ${id}`, exact: true });
  await link.focus();
  await page.keyboard.press('Enter');
  await expect(details.getByText(id, { exact: true })).toBeVisible();
  expect(mutations).toEqual([]);
  expect(errors).toEqual([]);
});

test('direct run lookup distinguishes loading, unknown identifiers, and API failure with recovery', async ({ page, request }) => {
  await page.clock.install();
  let release!: () => void;
  const gate = new Promise<void>((resolve) => { release = resolve; });
  await page.route('**/api/runs', async (route) => { await gate; await route.continue(); });
  try {
    await page.goto('/runs/unknown-run');
    await expect(page.getByText('Loading run…', { exact: true })).toBeVisible();
    await expect(page.getByRole('heading', { name: 'Run not found', exact: true })).toHaveCount(0);
    release();
    await expect(page.getByRole('heading', { name: 'Run not found', exact: true })).toBeVisible();
  } finally {
    release();
    await page.unroute('**/api/runs');
  }
  await addRun('recovered-run', 'Completed', 1794225600, 'cicd');
  service('stop', 'api');
  try {
    await page.goto('/runs/recovered-run');
    await expect(page.getByRole('heading', { name: 'Unable to load run', exact: true })).toBeVisible();
    await expect(page.getByRole('alert')).toContainText('unavailable');
    await expect(page.getByRole('heading', { name: 'Run not found', exact: true })).toHaveCount(0);
    service('start', 'api');
    await expect.poll(async () => (await request.get('/api/runs')).status()).toBe(200);
    await page.getByRole('button', { name: 'Try again', exact: true }).click();
    await expect(page.getByRole('region', { name: 'Run details', exact: true })).toContainText('recovered-run');
    await expect(page.getByRole('alert')).toHaveCount(0);
  } finally {
    service('start', 'api');
  }
});

test('selected run updates its lifecycle and elapsed time automatically, then keeps a recorded duration fixed', async ({ page, request }) => {
  await addRun('live-run', 'Pending', 1794225600, 'cicd');
  await page.clock.install({ time: new Date('2026-11-09T12:01:05Z') });
  await page.goto('/runs/live-run');
  const details = page.getByRole('region', { name: 'Run details', exact: true });
  await expect(details.getByText('Pending', { exact: true })).toBeVisible();
  await expect(details.getByText('Unavailable', { exact: true })).toHaveCount(3);
  const update = async (data: object) => {
    const response = await request.put(
      `http://127.0.0.1:${process.env.API_TEST_PORT}/app/cicd/api/runs/live-run`, { data },
    );
    expect(response.ok()).toBe(true);
  };
  await update({ status: 'Running', time_started: 1794225600 });
  await page.clock.runFor(5000);
  await expect(details.getByText('Running', { exact: true })).toBeVisible();
  await expect(details.getByText('1m 10s', { exact: true })).toBeVisible();
  await page.clock.runFor(1000);
  await expect(details.getByText('1m 11s', { exact: true })).toBeVisible();
  // A recorded end fixes execution duration even if status has not caught up yet.
  await update({ status: 'Running', time_ended: 1794225675 });
  await page.clock.runFor(4000);
  await expect(details.getByText('9 Nov 2026, 12:01:15 UTC', { exact: true })).toBeVisible();
  await page.clock.runFor(2000);
  await expect(details.getByText('1m 15s', { exact: true })).toBeVisible();
  await update({ status: 'Completed' });
  await page.clock.runFor(5000);
  await expect(details.getByText('Completed', { exact: true })).toBeVisible();
  await page.clock.runFor(10_000);
  await expect(details.getByText('1m 15s', { exact: true })).toBeVisible();
  await expect(page).toHaveURL(/\/runs\/live-run$/);
});

test('operator sees every run lifecycle, honest missing times, and recorded or absent failure detail', async ({ page }) => {
  const reason = 'Build failed: καλημέρα <img src=x onerror=alert(1)>\nCompiler exited with code 2.';
  for (const status of ['Pending', 'Running', 'Completed', 'Failed', 'Canceled', 'Paused']) {
    await addRun(`run-${status}`, status, 0, 'cicd', status === 'Failed' ? { error: reason } : {});
    await page.goto(`/runs/run-${status}`);
    const details = page.getByRole('region', { name: 'Run details', exact: true });
    await expect(details.getByText(status === 'Paused' ? 'Unknown (Paused)' : status, { exact: true })).toBeVisible();
    await expect(details.getByText('Unavailable', { exact: true })).toHaveCount(4);
    await expect(details.getByText(/1970/)).toHaveCount(0);
    if (status === 'Failed') {
      await expect(details.getByRole('heading', { name: 'Failure reason', exact: true })).toBeVisible();
      await expect(details.getByText(reason, { exact: true })).toBeVisible();
      await expect(page.locator('img')).toHaveCount(0);
    } else {
      await expect(details.getByRole('heading', { name: 'Failure reason', exact: true })).toHaveCount(0);
    }
  }
  await addRun('no-reason', 'Failed', 1794225600, 'cicd', {
    time_started: null, time_ended: null, error: null,
  });
  await page.goto('/runs/no-reason');
  await expect(page.getByRole('region', { name: 'Run details', exact: true }))
    .toContainText('No failure detail was recorded for this run.');
});
