# Recent-runs dashboard

The dashboard is a separate TypeScript/Vite module and Nginx container. It shows
stored runs across repositories, newest-triggered first (equal timestamps use
ascending run ID). Refresh reads the existing API; the frontend stores no history.
The dashboard is read-only. Open a run to inspect its lifecycle and recorded
failure alongside saved logs, including completed-run history, and follows newly
persisted output. Viewer search/download controls follow in ticket #118.

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

## Saved logs

Opening a known run reads its persisted records from the log collector. Each
record retains its ID, timestamp string, numeric severity (shown as INFO, WARN,
ERROR or DEBUG), and literal message text. Equal messages in different records
remain separate lines. The keyboard-focusable pane wraps long lines and contains
scrolling within the page. Markup-like output and Unicode are displayed as text.

The UI sorts each snapshot by the stored timestamp string ascending,
then record ID ascending. The backend read contract stays unchanged and does not
guarantee response ordering. This is a deterministic presentation order, including
for older records; it does **not** reconstruct emission order. Timestamps are
stored strings, may use different formats, and can repeat. The UUID tie-breaker
is not a sequence or ingestion cursor. No schema, response envelope, record fields,
or WebSocket message format changes are required. Empty reads retain the existing
successful response with an empty array (the browser also accepts null data).
The UI proxy preserves escaped IDs, subject to the existing collector's route
constraints. IDs containing slashes cannot be read by that backend; a failed log
read leaves the run details usable.

The UI renders the numeric severity actually stored on each record, using
INFO=0, WARN=1, ERROR=2, DEBUG=3. The existing collector's constant declaration
stores every incoming severity as zero, so those records display INFO. The UI
does not reconstruct missing severity information or change backend ingestion.

Loading, no saved output, waiting for output on Pending/Running runs, and retrieval
failure have distinct messages. **Refresh saved logs** reloads the snapshot and
also retries failed reads. Requests time out after ten seconds and cannot overlap.
A failed refresh retains the previous output with an explicit stale warning;
even a previously empty snapshot remains stale rather than claiming a fresh empty
result. Successful retry clears that warning. Run details and their API refresh
continue independently during a log-collector outage. Navigating away aborts the
old log request, closes its subscription, cancels retries, clears its output, and
prevents late results from crossing runs.

## Following live output

The browser opens the collector's existing text-only WebSocket for the selected
run. Frames are update signals, never appended as unidentified records. Opening
the connection and receiving a frame each request a saved-record reconciliation.
A signal during a pending read queues another read after it completes, covering
initial-load/subscription overlap. Records are reconciled by persisted ID; equal
message text in distinct records remains separate. Existing row nodes are retained
when their content is unchanged.

Saved records are also read five seconds after each successful reconciliation,
including for completed, failed, and canceled runs. A terminal status change
requests an immediate reconciliation, but never ends periodic reads: asynchronously
ingested late logs and missed live notifications remain discoverable. The UI
recovers persisted records only; it cannot recover events the backend never saved.

The connection reports connecting, open, reconnecting, and unavailable states.
After three failed connection attempts it reports unavailable while continuing
to retry. Connection attempts time out after ten seconds; reconnects use delays
of 1, 2, 4, 8, 16, then at most 30 seconds, reset after a successful connection.
Failed saved-record reads retry with the same bounded backoff, independently of
the WebSocket. A connected socket does not clear a failed-read warning; only a
successful saved-record read does. Run API failures remain separate from both.

**Auto-scroll** is a keyboard-accessible checkbox, enabled initially. Turning it
off pauses viewport movement while ingestion continues. Existing text selection,
keyboard focus, and the visible reading position survive background updates,
including late records inserted earlier in the deterministic order. Turning it
back on immediately follows the loaded tail. Its setting survives background
run/status refreshes and navigation between runs.

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
Compose network as the existing API and log collector, and reuses backend history.

Use your existing backend `.env.docker` configuration. The UI requires no new
environment file or credentials. From the repository root:

```sh
make deploy-local-ui
```

This starts the API, log collector and their existing dependencies through the
per-service targets,
then builds and starts the UI. Open <http://127.0.0.1:3000>. `make deploy-local`
now includes the UI alongside the API, runner, and log collector.
`make deploy-dashboard` starts the API, log collector, and UI through those same
per-service targets. `docker/docker-compose.dev-server.yaml` also includes the
UI on its existing `cicd_net` network.

The UI has these overridable defaults:

| Variable | Container default | Development server default |
| --- | --- | --- |
| `API_UPSTREAM` | `http://api:8080` | `http://127.0.0.1:8080` |
| `LOGCOLLECTOR_UPSTREAM` | `http://logcollector:8081` | `http://127.0.0.1:8081` |
| `UI_PORT` | `3000` | Not used |
| `UI_DEV_PORT` | Not used | `5173` |
| `TAG` | `latest` | Not used |

For example, `UI_PORT=3001 make deploy-local-ui` changes the container's localhost
port. For a different backend, set `API_UPSTREAM` to its reachable origin
(scheme, hostname, and port, without a path or trailing slash). Nginx substitutes
this at container startup; rebuilding the browser bundle is unnecessary.
Set `LOGCOLLECTOR_UPSTREAM` the same way for a different log collector origin.

The browser requests `GET /api/runs` from the UI origin. Nginx forwards that read
to `GET /app/cicd/api/runs` on the API. It resolves container DNS at request time.
The browser also reads `GET /api/logs/<encoded-run-id>`, forwarded to
`GET /app/cicd/logs/<encoded-run-id>` on the log collector. Live connections use
`/api/logs/ws/<encoded-run-id>`, forwarded to `/app/cicd/logs/ws/<encoded-run-id>`
with HTTP/1.1 Upgrade/Connection headers and the browser's original Host/Origin.
This preserves the collector's existing same-origin check without a CORS change.
Both Nginx and Vite route HTTP and WebSocket traffic to `LOGCOLLECTOR_UPSTREAM`.
Any gateway in front of the UI must forward WebSocket upgrades and preserve the
browser Host. Nginx closes an idle live stream after 60 seconds because the
existing protocol has no heartbeat; the UI reconnects and reconciles automatically.
Only these read routes are proxied; other API routes and mutations are rejected. Browser configuration contains no GitHub,
database, or RabbitMQ credentials, and no CORS changes are required.

After pulling an update, run `make deploy-local-ui` again to rebuild the UI.
Stop just the UI with `docker compose -f docker/docker-compose.ui.yaml down`.
Restarting or replacing it preserves backend history; the UI has no database or
data volume. Its Compose healthcheck checks the static server, while the page
shows API and log-collector connection failures independently.

## Development and image publication

Node 22.12+ (Node 24 in the production build) and npm are required locally:

```sh
npm ci --ignore-scripts --prefix ui
make dev-ui
make typecheck-ui
make build-ui
```

`make dev-ui` starts the existing local API and log collector before running Vite with
hot reload at <http://127.0.0.1:5173>. It can run alongside the container on port
3000. Set `UI_DEV_PORT`, `API_UPSTREAM`, or `LOGCOLLECTOR_UPSTREAM` to override its defaults. It fails
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
publisher sends logs through RabbitMQ to the real collector.
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

Saved-log journeys additionally cover collector-persisted RabbitMQ events, repeated
message identity, deterministic UI ordering and stored severity, historical rows,
active and completed empty states, loading, stale/initial failures and retry,
collector and UI restarts, delayed responses across selections, escaped run IDs,
literal Unicode/markup output, long lines, keyboard scrolling, and narrow layouts.
Live-log journeys additionally cover real WebSocket frames from broker-published
records, initial-load/subscription overlap, interrupted connections and catch-up,
late terminal output, run-switch isolation, and keyboard auto-scroll control with
stable reading position and text selection.

`make test-ui-dev` runs the run ordering, loading/empty, outage/recovery, and saved-log history journeys
through the actual Vite development server and the same disposable backend.

To run a single journey: `npm --prefix ui run test:browser -- --grep 'newest first'`.
Set `COMPOSE_BIN` to a standalone Compose executable if `docker compose` is not
on your PATH. Test traces/screenshots are under `ui/test-results/` on failure.
Validation runs through the Make targets locally; there is no UI GitHub Actions
workflow. Typechecking includes the browser harness.
