import { spawnSync } from 'node:child_process';
import pg from 'pg';
import amqp from 'amqplib';
import { test as base } from '@playwright/test';

if (!process.env.PGPORT || !process.env.POSTGRES_PASSWORD || !process.env.TEST_COMPOSE_ARGS || !process.env.UI_URL) {
  throw new Error('Use npm run test:browser to create the isolated fixture stack first.');
}

export const database = new pg.Pool({
  host: '127.0.0.1', port: Number(process.env.PGPORT),
  user: 'cicd', password: process.env.POSTGRES_PASSWORD, database: 'cicd',
});

// Files share this pool when Playwright reuses a worker; close it only at worker teardown.
export const test = base.extend<{}, { databaseConnection: void }>({
  databaseConnection: [async ({}, use) => {
    try {
      await use();
    } finally {
      await database.end();
    }
  }, { scope: 'worker', auto: true }],
});

export async function resetRuns() {
  await database.query('TRUNCATE api.pipeline_runs');
  await database.query('TRUNCATE logcollector.log_entries');
}

export async function addRun(id: string, status: string, time: number, repository = 'recipe-book', overrides = {}) {
  const run = {
    id, status, repository, branch: 'main', error: '', trigger_type: 'Push',
    triggered_by: 'SpyrosMoux', time_triggered: time, time_started: 0, time_ended: 0,
    ...overrides,
  };
  const fields = Object.keys(run);
  const placeholders = fields.map((_, i) => '$' + (i + 1)).join(',');
  await database.query(
    `INSERT INTO api.pipeline_runs (${fields.join(',')}) VALUES (${placeholders})`,
    Object.values(run),
  );
}

export function service(action: 'stop' | 'start' | 'restart', name: 'api' | 'ui' | 'logcollector') {
  const result = spawnSync(process.env.TEST_COMPOSE_BIN!, [
    ...JSON.parse(process.env.TEST_COMPOSE_ARGS!), action, name,
  ], { encoding: 'utf8' });
  if (result.status !== 0) throw new Error(result.stderr);
}

// Shared fixture publisher for subsequent log journeys: the real collector persists each record.
export async function publishLog(runId: string, message: string, overrides = {}) {
  const connection = await amqp.connect(process.env.RABBITMQ_URL!);
  try {
    const channel = await connection.createConfirmChannel();
    await channel.assertQueue('logs', { durable: true });
    channel.sendToQueue('logs', Buffer.from(JSON.stringify({
      runId, timestamp: '2026-11-09T12:00:00Z', logLevel: 'INFO', message, ...overrides,
    })), { persistent: true });
    await channel.waitForConfirms();
    await channel.close();
  } finally {
    await connection.close();
  }
}
