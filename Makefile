-include .env

.PHONY: test
test:
	@go test -v ./...
