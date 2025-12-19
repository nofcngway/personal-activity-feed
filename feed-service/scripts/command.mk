.PHONY: generate-api
generate-api:
	@./scripts/generate.sh

.PHONY: mock
mock:
	@mockery

.PHONY: run
run:
	@configPath=./config.yaml go run ./cmd/app


