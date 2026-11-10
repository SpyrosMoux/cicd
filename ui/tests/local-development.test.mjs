import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import { fileURLToPath } from 'node:url';
import test from 'node:test';

test('the standalone UI configures with defaults and reuses existing backend services', () => {
  const env = { ...process.env };
  for (const key of ['POSTGRES_USER', 'POSTGRES_PASSWORD', 'POSTGRES_DB',
    'RABBITMQ_USER', 'RABBITMQ_PASSWORD', 'GH_TOKEN', 'GH_WEBHOOK_SECRET']) delete env[key];
  // Empty overrides must also use Compose's default values.
  env.API_UPSTREAM = '';
  env.UI_PORT = '';
  env.TAG = '';
  const executable = env.COMPOSE_BIN || 'docker';
  const prefix = env.COMPOSE_BIN ? [] : ['compose'];
  const result = spawnSync(executable, [...prefix, '--env-file', '/dev/null',
    '-f', 'docker/docker-compose.ui.yaml', 'config', '--format', 'json'], {
    cwd: fileURLToPath(new URL('../../', import.meta.url)), env, encoding: 'utf8',
  });
  assert.ifError(result.error);
  assert.equal(result.status, 0, result.stderr);
  assert.deepEqual(Object.keys(JSON.parse(result.stdout).services), ['ui']);
});
