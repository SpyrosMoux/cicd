# FlowForge

FlowForge is a self-hostable CI/CD platform that runs pipelines using Docker containers. A separately deployed, read-only dashboard shows recent pipeline runs. The platform can be deployed using Docker Compose or Kubernetes.

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
validation commands, see the [dashboard guide](docs/dashboard.md).

1. Clone the repository:
   ```sh
   git clone https://github.com/SpyrosMoux/cicd.git
   cd cicd
   ```
2. Configure the backend's `.env.docker` using the [setup guide](docs/setup.mdx).
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

- Docker (No additional dependencies required)

## Contributing

We welcome contributions! You can:
- Report issues in the repository.
- Open pull requests with improvements and new features.

## License

[GNU GENERAL PUBLIC LICENSE](LICENSE)

---
