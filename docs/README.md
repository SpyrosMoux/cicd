# FlowForge documentation

FlowForge's documentation is authored and published with the existing Mintlify
site. Pages live in this directory as MDX files; `docs.json` controls navigation.
The published entry point is [FlowForge Docs](https://docs.flowforge.spyrosmoux.com/).

## Pages

| Topic | Source |
| --- | --- |
| Introduction | [Introduction](introduction.mdx) |
| Architecture and service boundaries | [Architecture](architecture.mdx) |
| Models, tables, and message contracts | [Data models](data-model.mdx) |
| Webhook-to-dashboard walkthrough | [Execution flow](execution-flow.mdx) |
| Pipeline syntax and execution rules | [Pipeline reference](pipelines.mdx) |
| Domain vocabulary | [Glossary](glossary.mdx) |
| Contributing changes | [Contributing](contributing.mdx) |
| Local setup and troubleshooting | [Local development](development.mdx) |
| Dashboard deployment and behavior | [Dashboard](dashboard.mdx) |
| Server installation | [Installation](setup.mdx) |
| GitHub webhooks | [Webhook setup](essentials/setup-webhook.mdx) |
| Kubernetes | [Kubernetes setup](kubernetes-setup.mdx) |

## Preview and validate

From the repository root, use the installed Mintlify CLI:

```sh
cd docs
mint broken-links
mint validate --no-build
mint dev --no-open --port 3002
```

Port 3002 avoids the dashboard's default port 3000. Stop the preview with Ctrl+C.
Update this installation's `docs.json` configuration when adding pages. The
`--no-build` validation checks content and configuration without downloading or
building the preview client. Preview changed pages and diagrams, check internal
links, and run `git diff --check` before opening a PR.

## Editing conventions

Give every page a frontmatter title and description. Use extensionless site routes
such as `/architecture` for links between pages. Link to GitHub for source files
outside the Mintlify documentation root. Add pages to the appropriate navigation
group in `docs.json`.

Keep explanations in their owning guide and link to them elsewhere. The root
[CONTRIBUTING.md](https://github.com/SpyrosMoux/cicd/blob/main/CONTRIBUTING.md) directs GitHub contributors to the Mintlify
guide. The root [GLOSSARY.md](https://github.com/SpyrosMoux/cicd/blob/main/GLOSSARY.md) remains the canonical domain
vocabulary for the repository; keep the site's glossary synchronized when terms
change. Record newly agreed architectural decisions under `docs/adr/` following
the [domain documentation rules](agents/domain.md).

Never include credentials or rendered configuration containing secrets in docs.
