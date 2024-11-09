# TODO(spyrosmoux) accommodate repo merge
run-local-deps:
	docker compose -f docker-compose.deps.yaml up -d

run-local-runner: run-local-deps
	docker compose -f docker-compose.runner.yaml up -d

run-local-api: run-local-deps
	docker compose -f docker-compose.api.yaml up -d

run-local: run-local-deps run-local-runner run-local-api

proxy-webhook: # for local use only
	smee -u https://smee.io/vgX1mcHUonHXl1Hh -p 8080 -P /webhook
