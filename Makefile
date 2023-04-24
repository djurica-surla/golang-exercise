run:
	go run cmd/server/main.go

get-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.50.1

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix