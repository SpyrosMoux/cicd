# TODO(spyrosmoux) accommodate repo merge
run-local-deps:
	docker compose -f docker-compose.deps.yaml up -d

proxy-webhook: # for local use only
	smee -u https://smee.io/vgX1mcHUonHXl1Hh -p 8080 -P /webhook
