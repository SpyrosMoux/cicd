# FlowForge glossary

FlowForge receives repository events, executes matching pipelines, and exposes
their recorded progress and output. These terms describe that single domain.

## Language

**Pipeline**:
A repository-owned definition of triggers, variables, and ordered jobs.
_Avoid_: Workflow, when referring to a FlowForge pipeline.

**Pipeline run**:
One execution attempt of a matching pipeline, identified by a unique run ID.
_Avoid_: Job, when referring to the entire execution attempt.

**Trigger**:
A condition that selects a pipeline for a repository event, using the pushed
branch or the destination branch of a pull request.

**Job**:
A named group of steps inside a pipeline.
_Avoid_: Run, when referring to a group in the pipeline definition.

**Step**:
A named shell command or shell script within a job.

**Runner**:
The worker that prepares repository sources and executes a pipeline run.

**Log record**:
One persisted unit of run output with its own identity, timestamp, severity,
and message. Identical messages can belong to different records.

**Log collector**:
The service that stores run output and makes saved records and live update
notifications available to readers.

**Dashboard**:
The read-only interface for inspecting recorded pipeline runs and their output.

**Log snapshot**:
The saved records available to a reader from one successful retrieval. A snapshot
does not establish that all output for a run has arrived.
