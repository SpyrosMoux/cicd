# FlowForge

FlowForge is a self-hostable CI/CD platform with containerized runners that execute pipeline shell commands. A separately deployed, read-only dashboard shows recent pipeline runs. The platform includes Docker Compose and Kubernetes deployment configurations.

## Documentation

Start with the [Mintlify documentation](https://docs.flowforge.spyrosmoux.com/), or jump to:

- [Contributing](CONTRIBUTING.md) and [local development](https://docs.flowforge.spyrosmoux.com/development)
- [Architecture and service boundaries](https://docs.flowforge.spyrosmoux.com/architecture)
- [Data models and contracts](https://docs.flowforge.spyrosmoux.com/data-model)
- [From webhook to visible output](https://docs.flowforge.spyrosmoux.com/execution-flow)
- [Pipeline syntax and execution rules](https://docs.flowforge.spyrosmoux.com/pipelines)
- [Domain glossary](GLOSSARY.md)

## How It Works

For full documentation, visit: [FlowForge Docs](https://docs.flowforge.spyrosmoux.com/)

1. You set up FlowForge by providing a URL for its API.
2. You configure webhooks on your repositories to send events to this API.
3. When an event is received, FlowForge determines which pipeline to run.
4. The event is published to RabbitMQ.
5. A runner picks up the pipeline from the queue and executes it using Docker.

## Features

- Runs pipelines inside Docker containers.
- Can be deployed in Kubernetes.
- Support for autoscaling in Kubernetes.
- Lightweight and easy to set up with Docker Compose.
- Standalone recent-runs dashboard with a production image and isolated browser validation.

## Installation

For the standalone dashboard, local deployment, and development and
validation commands, see the [dashboard guide](https://docs.flowforge.spyrosmoux.com/dashboard).

1. Clone the repository:
   ```sh
   git clone https://github.com/SpyrosMoux/cicd.git
   cd cicd
   ```
2. Configure the local environment and start dependencies using the [development guide](https://docs.flowforge.spyrosmoux.com/development#configure-the-local-backend).
3. Start the platform using Docker Compose:
   ```sh
   make deploy-local
   ```
4. Set up webhooks on your repositories to send events to your FlowForge API URL.

## Example Workflow

1. A developer pushes a commit to GitHub.
2. GitHub sends a `push` event webhook to FlowForge.
3. The FlowForge API processes the event and publishes a job to RabbitMQ.
4. A runner picks up the job from the queue and runs the pipeline inside a Docker container.

## Requirements

- Docker Engine and Docker Compose for container deployments; Make for the commands above.
- Go and Node/npm for source development; see the [contributing guide](https://docs.flowforge.spyrosmoux.com/contributing#getting-started).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for setup, coding conventions, validation,
and the issue/PR process. Bug reports and proposals belong in
[GitHub Issues](https://github.com/SpyrosMoux/cicd/issues).

## License

[GNU GENERAL PUBLIC LICENSE](LICENSE)

---
