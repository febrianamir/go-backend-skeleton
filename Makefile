run-api:
	@air -c .air/.air.api.toml

run-worker:
	@air -c .air/.air.worker.toml

run-scheduler:
	@go run apps/scheduler/main.go

generate-hmac-key:
	@openssl rand 32 | base64 | tr '+/' '-_' | tr -d '='

