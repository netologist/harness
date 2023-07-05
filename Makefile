prerequisites:
	@go install -mod=mod -v \
		github.com/golangci/golangci-lint/cmd/golangci-lint \
		github.com/golang/mock/mockgen

install:
	@go mod download
	@go mod tidy

test:
	@go test -v --cover -count=1 -race ./...

lint:
	@golangci-lint run

generate-mocks:
	@go generate `go list ./... | grep -v tests/e2e`