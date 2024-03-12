lint: lint_deps
	@golangci-lint --version
	@golangci-lint run --fix | tee lint.err # https://golangci-lint.run/usage/install/#local-installation

lint_deps: gofmt vet

vet:
	@go vet ./...

gofmt:
	@GO111MODULE=off gofmt -l $(shell find . -type f -name '*.go'| grep -v "/vendor/\|/.git/")

test:
		LOGGER_PATH=nul CGO_CFLAGS=-Wno-undef-prefix go test -test.v -timeout 99999m -cover ./...
