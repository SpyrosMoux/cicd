# Recent-runs dashboard

The dashboard is a separate TypeScript/Vite module and Nginx container. It shows
stored runs across repositories, newest-triggered first (equal timestamps use
ascending run ID). Refresh reads the existing API; the frontend stores no history.
The dashboard is read-only. Open a run to inspect its lifecycle and recorded
failure. The log viewer follows in subsequent tickets under spec #112.

## Investigating a run

Select a run identifier in the list to open `/runs/<encoded-run-id>`. This is a
stable link: it works when pasted into a new tab, refreshed, or revisited using
browser back/forward. The packaged Nginx server and Vite development server both
serve the dashboard for direct run links. **Back to runs** returns to the list.
Available list filters, reading position, and the originating link's keyboard
focus are retained within that tab's navigation history, including a detail-page
refresh. A newly opened direct link starts with the default list context.

Details show the identifier, repository, branch, status, trigger type, triggering
actor, and UTC lifecycle times from the existing run-list API. Zero or missing
timestamps are unavailable. Duration uses actual execution timestamps: running
work advances locally every second; a recorded end fixes the duration. Unstarted
work has no duration. Failed runs show the recorded failure as plain text, or
explain that no detail was recorded.

The same automatic refresh keeps the selected run current without changing its
URL or selection. Loading, a successfully resolved unknown identifier, and an
API failure have separate states. During an outage, usable details remain
visible with an explicit stale warning; automatic retries and **Try again**
provide recovery. Navigation uses the current URL when a response arrives, so a
late response cannot select a previously viewed run. No new API endpoint or
mutation operation is introduced.

## Filtering and monitoring

Use the **Repository** and **Status** selectors independently or together.
Repositories come from the complete loaded history. The status selector always
offers Pending, Running, Completed, Failed, and Canceled; unexpected recorded
statuses appear as Unknown with their original value. **Clear filters** restores
the full list. Counts show matches against the loaded total when filters are
active. An empty installation says **No runs yet**; filters with no matches say
**No matching runs**.

The page automatically reads fresh run history five seconds after each request
finishes, including after failures. Requests have a ten-second timeout and never
overlap; **Refresh runs** also restarts that schedule. New runs appear newest
first, and changed statuses immediately enter or leave the filtered list.
Background updates preserve filters, keyboard focus, unchanged rows and selected
text, and the reading position when unchanged visible rows remain. Newly seen
repositories become selectable. An active selection stays available even if its
last run disappears from the latest snapshot.

If the API becomes unavailable, previously loaded history remains usable with a
warning that it may be stale. Filters still work against that retained snapshot.
The page retries automatically and clears the warning after a successful read;
**Try again** also allows immediate recovery. A failure before any successful
load shows **Unable to load runs**, never a successful empty result.

## Local deployment

The UI follows the backend's per-service Compose and Makefile pattern. Its
`docker/docker-compose.ui.yaml` defines only the UI, joins the same default
Compose network as the existing API, and reuses backend history.

Use your existing backend `.env.docker` configuration. The UI requires no new
environment file or credentials. From the repository root:

```sh
make deploy-local-ui
```

This starts the API and its existing dependencies through `deploy-local-api`,
then builds and starts the UI. Open <http://127.0.0.1:3000>. `make deploy-local`
now includes the UI alongside the API, runner, and log collector.
`make deploy-dashboard` starts the API, log collector, and UI through those same
per-service targets. `docker/docker-compose.dev-server.yaml` also includes the
UI on its existing `cicd_net` network.

The UI has these overridable defaults:

| Variable | Container default | Development server default |
| --- | --- | --- |
| `API_UPSTREAM` | `http://api:8080` | `http://127.0.0.1:8080` |
| `UI_PORT` | `3000` | Not used |
| `UI_DEV_PORT` | Not used | `5173` |
| `TAG` | `latest` | Not used |

For example, `UI_PORT=3001 make deploy-local-ui` changes the container's localhost
port. For a different backend, set `API_UPSTREAM` to its reachable origin
(scheme, hostname, and port, without a path or trailing slash). Nginx substitutes
this at container startup; rebuilding the browser bundle is unnecessary.

The browser requests `GET /api/runs` from the UI origin. Nginx forwards that read
to `GET /app/cicd/api/runs` on the API. It resolves container DNS at request time.
Only this read route is proxied. Browser configuration contains no GitHub,
database, or RabbitMQ credentials, and no CORS changes are required.

After pulling an update, run `make deploy-local-ui` again to rebuild the UI.
Stop just the UI with `docker compose -f docker/docker-compose.ui.yaml down`.
Restarting or replacing it preserves backend history; the UI has no database or
data volume. Its Compose healthcheck checks the static server, while the page
shows API connection failures.

## Development and image publication

Node 22.12+ (Node 24 in the production build) and npm are required locally:

```sh
npm ci --ignore-scripts --prefix ui
make dev-ui
make typecheck-ui
make build-ui
```

`make dev-ui` starts the existing local API dependencies before running Vite with
hot reload at <http://127.0.0.1:5173>. It can run alongside the container on port
3000. Set `UI_DEV_PORT` or `API_UPSTREAM` to override its defaults. It fails
clearly if the selected port is already occupied instead of choosing another.
Both settings remain server-side; do not put secrets in `VITE_*` variables.
Production uses a multi-stage image with fixed base-image versions and an npm
lockfile. The final Nginx image runs as UID 101 on container port 8080; Node and
test dependencies are absent. Compose still publishes localhost port 3000.

The existing registry/tag conventions also apply to the UI:

```sh
make build-ui-docker TAG=YOUR_VERSION PLATFORM=linux/amd64
```

This target builds **and pushes** `ghcr.io/spyrosmoux/cicd/ui:YOUR_VERSION`, like
the existing backend targets. Use `PLATFORM=linux/arm64` for an ARM mini PC. A
plain `docker build` as above builds locally without publishing. `build-all-docker`
now includes the UI. For published UI images, set `TAG`, then run `docker compose
-f docker/docker-compose.ui.yaml pull` followed by `up -d --no-build` with that
same Compose file.

## Isolated browser validation

With Docker Engine and Compose v2.24.4+ available:

```sh
npm ci --ignore-scripts --prefix ui
cd ui
npx playwright install --with-deps chromium
cd ..
make test-ui
```

The harness combines the same per-service Compose files with a test-only
override. It builds the production UI image and real API/log-collector images,
starts disposable PostgreSQL/RabbitMQ with freshly generated credentials, and uses
random localhost ports and a unique Compose project. It seeds deterministic
records in that database, exercises real HTTP reads through Nginx, and removes
only its own stack and volumes afterward. No live GitHub event, runner execution,
mocked API response, or production test endpoint is used. Its shared fixture
publisher sends logs through RabbitMQ to the real collector for later slices.
Journeys cover combined filters, no matches/reset, new repositories and runs,
status changes through the existing API, automatic outage recovery, and bounded
refresh during a slow real request. Browser clock control advances refresh time;
API responses are never substituted. Keyboard and desktop/narrow-screen checks
include the filters, and a reading-position journey preserves selected text as
new work arrives.

Run-detail journeys additionally cover direct links (including encoded IDs),
refresh/back/forward, retained filters/focus/reading position, navigation before
initial loading finishes, every lifecycle state, zero/null times, live duration
and fixed recorded ends, recorded/missing failure detail, unknown IDs, API outage
recovery with retained stale details, and selection changes during a delayed real
response. Desktop/narrow-screen and keyboard checks include metadata and failure
text containing Unicode, markup-like content, and long lines.

`make test-ui-dev` runs the ordering, loading/empty, and outage/recovery journeys
through the actual Vite development server and the same disposable backend.

To run a single journey: `npm --prefix ui run test:browser -- --grep 'newest first'`.
Set `COMPOSE_BIN` to a standalone Compose executable if `docker compose` is not
on your PATH. Test traces/screenshots are under `ui/test-results/` on failure.
Validation runs through the Make targets locally; there is no UI GitHub Actions
workflow. Typechecking includes the browser harness.
