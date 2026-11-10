import { element } from './dom';
import { formatDuration, formatTime, statusLabel, type Run } from './runs';

export function renderLookupState(state: 'loading' | 'unknown' | 'unavailable') {
  const [heading, description] = {
    loading: ['Loading run…', 'Looking up this run.'],
    unknown: ['Run not found', 'No run with this identifier is available.'],
    unavailable: ['Unable to load run', 'The run could not be looked up. Try again when the API is reachable.'],
  }[state];
  const content = element('div', 'state');
  content.append(element('h2', '', heading), element('p', '', description));
  return content;
}

export function renderDetails(run: Run) {
  const content = element('div');
  const fields = element('dl', 'detail-fields');
  const values = [
    ['Run identifier', run.id], ['Repository', run.repository], ['Branch', run.branch],
    ['Status', statusLabel(run.status)], ['Trigger type', run.trigger_type],
    ['Triggered by', run.triggered_by], ['Triggered', formatTime(run.time_triggered)],
    ['Started', formatTime(run.time_started)], ['Ended', formatTime(run.time_ended)],
    ['Duration', formatDuration(run)],
  ];
  for (const [label, value] of values) {
    const field = element('div');
    const content = element('dd', '', value || 'Unavailable');
    if (label === 'Duration') content.dataset.duration = '';
    field.append(element('dt', '', label), content);
    fields.append(field);
  }
  content.append(fields);
  if (run.status === 'Failed') {
    const failure = element('section', 'failure-detail');
    failure.append(element('h2', '', 'Failure reason'),
      element('p', '', run.error?.trim() ? run.error : 'No failure detail was recorded for this run.'));
    content.append(failure);
  }
  return content;
}
