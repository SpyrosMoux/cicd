build:
	go mod tidy
	go build -o ${PWD}/bin/api cmd/main.go

docker-build:
	docker build -t ghcr.io/spyrosmoux/api .

docker-push: docker-build
	docker push ghcr.io/spyrosmoux/api

run-local-deps:
	docker compose -f docker-compose.deps.yaml up -d

run-local: build run-local-deps
	go run cmd/main.go


run-docker: run-local-deps
	docker compose -f docker-compose.yaml up -d

proxy-webhook: # for local use only
	smee -u https://smee.io/vgX1mcHUonHXl1Hh -p 8080 -P /webhook
