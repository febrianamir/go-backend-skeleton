run-api:
	@air -c .air/.air.api.toml

run-worker:
	@air -c .air/.air.worker.toml

generate-hmac-key:
	@openssl rand 32 | base64 | tr '+/' '-_' | tr -d '='
