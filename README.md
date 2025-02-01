# TODO

- Create a new runs package
- When adding a new pipeline, add validation, so each saved pipeline has the correct structure

- When an event occurs
  - Get the saved pipelines from the db based on the event's repo
  - Based on the event's branch, checkout each pipeline's yaml
  - Check each pipeline's trigger
  - Create a new run for each triggered pipeline

# > [!NOTE]
>
> For now I cannot think of a way to check if an invalid yaml is triggered. This would help
> by showing the user that although this pipeline is triggered, it cannot run because the yaml
> is invalid
