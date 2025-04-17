# Define variables for reuse
GO_CMD = go build -o
API_DIR = cmd/api
RUNNER_DIR = cmd/runner
LOGCOLLECTOR_DIR = cmd/logcollector
BUILD_DIR = build
DOCKER_REPO = ghcr.io/spyrosmoux/cicd
SMEE_URL = https://smee.io/vgX1mcHUonHXl1Hh
PROXY_PORT ?= 8080
PLATFORM = linux/amd64
TAG ?= latest

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

build-logcollector:
	$(GO_CMD) $(BUILD_DIR)/logcollector/ $(LOGCOLLECTOR_DIR)/main.go

build-all: build-api build-runner build-logcollector

# Docker builds with optional tag argument
build-api-docker:
	docker build -t $(DOCKER_REPO)/api:$(TAG) . -f docker/Dockerfile.api --platform $(PLATFORM)
	docker push $(DOCKER_REPO)/api:$(TAG)

build-runner-docker:
	docker build -t $(DOCKER_REPO)/runner:$(TAG) . -f docker/Dockerfile.runner --platform $(PLATFORM)
	docker push $(DOCKER_REPO)/runner:$(TAG)

build-logcollector-docker:
	docker build -t $(DOCKER_REPO)/logcollector:$(TAG) . -f docker/Dockerfile.logcollector --platform $(PLATFORM)
	docker push $(DOCKER_REPO)/logcollector:$(TAG)

build-all-docker: build-api-docker build-runner-docker build-logcollector-docker

# Local deployments
deploy-local-deps:
	docker compose -f docker/docker-compose.deps.yaml up -d

deploy-local-runner: deploy-local-deps
	docker compose -f docker/docker-compose.runner.yaml up -d

deploy-local-api: deploy-local-deps
	docker compose -f docker/docker-compose.api.yaml up -d

deploy-local-logcollector: deploy-local-deps
	docker compose -f docker/docker-compose.logcollector.yaml up -d

deploy-local: deploy-local-deps deploy-local-runner deploy-local-api deploy-local-logcollector

deploy-dev-server:
	docker compose -f docker/docker-compose.dev-server.yaml pull
	docker compose -f docker/docker-compose.dev-server.yaml up -d

# Proxy webhook for local development
proxy-webhook:
	smee -u $(SMEE_URL) -p $(PROXY_PORT) -P /app/cicd/api/gh/webhook

# Clean up builds and Docker containers
clean:
	rm -rf $(BUILD_DIR)/*
	docker compose -f docker/docker-compose.deps.yaml down
	docker compose -f docker/docker-compose.runner.yaml down
	docker compose -f docker/docker-compose.api.yaml down
	docker compose -f docker/docker-compose.logcollector.yaml down

k8s:
	kubectl apply -k kubernetes/kustomize/overlays/flowforge

k8s-down:
	kubectl delete -k kubernetes/kustomize/overlays/flowforge

install-keda:
	kustomize build kubernetes/kustomize/overlays/keda --enable-helm > keda.yaml
	kubectl apply --server-side -f keda.yaml
