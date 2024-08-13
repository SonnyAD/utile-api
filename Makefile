docs: 
	swag fmt && swag init

lint: 
	golangci-lint run

.PHONY: docs