# Define variables for reuse
GO_CMD = go build -o
API_DIR = cmd/api
RUNNER_DIR = cmd/runner
BUILD_DIR = build
DOCKER_REPO = ghcr.io/spyrosmoux/cicd
SMEE_URL = https://smee.io/vgX1mcHUonHXl1Hh

# Default targets
.PHONY: build-api build-runner build-all \
        build-api-docker build-runner-docker build-all-docker \
        deploy-local-deps deploy-local-runner deploy-local-api deploy-local \
        proxy-webhook clean

# Binary builds
build-api:
	$(GO_CMD) $(BUILD_DIR)/api/ $(API_DIR)/main.go

build-runner:
	$(GO_CMD) $(BUILD_DIR)/runner/ $(RUNNER_DIR)/main.go

build-all: build-api build-runner

# Docker builds with optional tag argument
build-api-docker:
	docker build -t $(DOCKER_REPO)/api:${TAG:=latest} . -f docker/Dockerfile.api
	docker push $(DOCKER_REPO)/api:${TAG:=latest}

build-runner-docker:
	docker build -t $(DOCKER_REPO)/runner:${TAG:=latest} . -f docker/Dockerfile.runner
	docker push $(DOCKER_REPO)/runner:${TAG:=latest}

build-all-docker: build-api-docker build-runner-docker

# Local deployments
deploy-local-deps:
	docker compose -f docker/docker-compose.deps.yaml up -d

deploy-local-runner: deploy-local-deps
	docker compose -f docker/docker-compose.runner.yaml up -d

deploy-local-api: deploy-local-deps
	docker compose -f docker/docker-compose.api.yaml up -d

deploy-local: deploy-local-deps deploy-local-runner deploy-local-api

# Proxy webhook for local development
proxy-webhook:
	smee -u $(SMEE_URL) -p 8080 -P /api/webhook

# Clean up builds and Docker containers
clean:
	rm -rf $(BUILD_DIR)/*
	docker compose -f docker/docker-compose.deps.yaml down
	docker compose -f docker/docker-compose.runner.yaml down
	docker compose -f docker/docker-compose.api.yaml down
