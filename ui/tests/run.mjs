import { spawn, spawnSync } from 'node:child_process';
import { randomUUID } from 'node:crypto';
import { createServer } from 'node:net';
import { fileURLToPath } from 'node:url';

const cwd = fileURLToPath(new URL('../', import.meta.url));
const development = process.argv.includes('--dev');
const browserArgs = process.argv.slice(2).filter((arg) => arg !== '--dev');
const project = `flowforge-browser-${process.pid}-${randomUUID().slice(0, 8)}`;
// Keep the UI origin and Vite's API target stable across container restarts.
const portProbes = [createServer(), createServer(), createServer()];
await Promise.all(portProbes.map((probe) =>
  new Promise((resolve, reject) => {
    probe.once('error', reject);
    probe.listen(0, '127.0.0.1', resolve);
  })));
const [uiPort, apiPort, logcollectorPort] = portProbes.map((probe) => probe.address().port);
await Promise.all(portProbes.map((probe) => new Promise((resolve) => probe.close(resolve))));
const env = {
  ...process.env,
  COMPOSE_PROJECT_NAME: project,
  TAG: project,
  UI_TEST_PORT: String(uiPort),
  API_TEST_PORT: String(apiPort),
  LOGCOLLECTOR_TEST_PORT: String(logcollectorPort),
  POSTGRES_PASSWORD: randomUUID(),
  RABBITMQ_PASSWORD: randomUUID(),
  GH_WEBHOOK_SECRET: randomUUID(),
  GH_TOKEN: randomUUID(),
  // Deliberately override the default: the packaged journey proves deploy-time routing.
  API_UPSTREAM: 'http://fixture-api:8080', // NOSONAR: Plain HTTP is confined to this disposable Docker network; the fixture API has no TLS listener.
  LOGCOLLECTOR_UPSTREAM: 'http://fixture-logs:8081', // NOSONAR: The fixture collector has no TLS listener on this disposable network.
};
const executable = process.env.COMPOSE_BIN || 'docker';
const prefix = process.env.COMPOSE_BIN ? [] : ['compose'];
const args = [...prefix, '-p', project,
  '-f', '../docker/docker-compose.deps.yaml',
  '-f', '../docker/docker-compose.api.yaml',
  '-f', '../docker/docker-compose.logcollector.yaml',
  '-f', '../docker/docker-compose.ui.yaml',
  '-f', 'tests/compose.yaml'];
let devServer;

function compose(command, capture = false) {
  const result = spawnSync(executable, [...args, ...command], {
    cwd, env, encoding: 'utf8', stdio: capture ? 'pipe' : 'inherit',
  });
  if (result.error || result.status !== 0) {
    throw new Error(result.error?.message ?? result.stderr ?? `Compose failed: ${command.join(' ')}`);
  }
  return result.stdout?.trim();
}

try {
  // Explicit services keep an operator's COMPOSE_PROFILES from starting a runner.
  compose(['up', '-d', '--build', '--wait', '--wait-timeout', '180', 'ui', 'logcollector']);
  env.UI_URL = `http://${compose(['port', 'ui', '8080'], true)}`;
  env.PGPORT = compose(['port', 'api-db', '5432'], true).split(':').at(-1);
  env.RABBITMQ_URL = `amqp://cicd:${encodeURIComponent(env.RABBITMQ_PASSWORD)}@${compose(['port', 'rabbitmq', '5672'], true)}`;
  env.LOGCOLLECTOR_URL = `http://${compose(['port', 'logcollector', '8081'], true)}`;
  env.TEST_COMPOSE_BIN = executable;
  env.TEST_COMPOSE_ARGS = JSON.stringify(args);
  if (development) {
    process.env.API_UPSTREAM = `http://${compose(['port', 'api', '8080'], true)}`;
    process.env.LOGCOLLECTOR_UPSTREAM = env.LOGCOLLECTOR_URL;
    process.env.UI_DEV_PORT = '0';
    const vite = await import('vite');
    devServer = await vite.createServer({ root: cwd, configFile: fileURLToPath(new URL('../vite.config.ts', import.meta.url)) });
    await devServer.listen();
    env.UI_URL = devServer.resolvedUrls.local[0];
  }
  process.exitCode = await new Promise((resolve, reject) => {
    const playwright = fileURLToPath(new URL('../node_modules/@playwright/test/cli.js', import.meta.url));
    const browser = spawn(process.execPath, [playwright, 'test', ...browserArgs], { cwd, env, stdio: 'inherit' });
    browser.once('error', reject);
    browser.once('exit', (code) => resolve(code ?? 1));
  });
  if (process.exitCode) compose(['logs', '--no-color', '--tail', '60']);
} catch (error) {
  console.error(error);
  process.exitCode = 1;
  compose(['logs', '--no-color', '--tail', '60']);
} finally {
  await devServer?.close();
  // This project name is generated here; no operator stack or volume is selected.
  compose(['down', '--volumes', '--remove-orphans']);
}
