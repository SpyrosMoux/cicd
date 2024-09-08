build:
	go mod tidy
	go build -o ${PWD}/bin/api cmd/main.go

docker-build:
	docker build -t ghcr.io/spyrosmoux/api .

run-local: build
	docker compose -f docker-compose.deps.yaml up -d
	go run cmd/main.go

run-docker: docker-build
	docker compose -f docker-compose.deps.yaml up -d
	docker compose -f docker-compose.yaml up -d

proxy-webhook: # for local use only
	smee -u https://smee.io/vgX1mcHUonHXl1Hh -p 8080 -P /webhook
