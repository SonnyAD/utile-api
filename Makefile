docs: 
	swag fmt && swag init

lint: 
	golangci-lint run

test:
	gotestsum -f short-verbose -- -short -coverprofile=cover.out ./...

.PHONY: docs