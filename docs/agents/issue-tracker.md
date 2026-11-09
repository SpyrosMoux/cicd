# Issue tracker: GitHub

Issues and specs live in SpyrosMoux/cicd on GitHub.
Use the gh CLI from this checkout, or pass --repo SpyrosMoux/cicd.

## Operations

- Create: gh issue create --title "..." --body-file <path>
- Read: gh issue view <number> --json number,title,body,labels,comments
- List: gh issue list --state open with appropriate label filters
- Comment: gh issue comment <number> --body-file <path>
- Labels: gh issue edit <number> --add-label "..." or --remove-label "..."
- Close: gh issue close <number> --comment "..."

Use native GitHub sub-issues for parent/child relationships and native
issue dependencies for blocking relationships. When unavailable, record
"Part of #<parent>" and "Blocked by: #<number>" in issue bodies.

When a skill says "publish to the issue tracker", create a GitHub issue.
When it says "fetch the relevant ticket", read the GitHub issue.

## Pull requests as a triage surface

**PRs as a request surface: no.**
